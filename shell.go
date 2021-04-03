package main

import (
	"fmt"
	"os"
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

)

func RunInput(input string) {
	// First try to run user input in Lua
	err := l.DoString(input)

	if err == nil {
		// If it succeeds, add to history and prompt again
		readline.AddHistory(input)
		readline.SaveHistory(homedir + "/.hilbish-history")
		bait.Em.Emit("command.success", nil)
		return
	}

	// Split up the input
	cmdArgs, cmdString := splitInput(input)
	// If there's actually no input, prompt again
	if len(cmdArgs) == 0 { return }

	// If alias was found, use command alias
	if aliases[cmdArgs[0]] != "" {
		cmdString = aliases[cmdArgs[0]] + strings.Trim(cmdString, cmdArgs[0])
		execCommand(cmdString)
		return
	}

	// If command is defined in Lua then run it
	if commands[cmdArgs[0]] {
		err := l.CallByParam(lua.P{
			Fn: l.GetField(
				l.GetTable(
					l.GetGlobal("commanding"),
					lua.LString("__commands")),
				cmdArgs[0]),
			NRet: 0,
			Protect: true,
		}, luar.New(l, cmdArgs[1:]))
		if err != nil {
			// TODO: dont panic
			panic(err)
		}
		readline.AddHistory(cmdString)
		readline.SaveHistory(homedir + "/.hilbish-history")
		return
	}

	// Last option: use sh interpreter
	switch cmdArgs[0] {
	case "exit":
		os.Exit(0)
	default:
		err := execCommand(cmdString)
		if err != nil {
			// If input is incomplete, start multiline prompting
			if syntax.IsIncomplete(err) {
				sb := &strings.Builder{}
				for {
					done := StartMultiline(cmdString, sb)
					if done {
						break
					}
				}
			} else {
				if code, ok := interp.IsExitStatus(err); ok {
					if code > 0 {
						bait.Em.Emit("command.fail", code)
					}
				}
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			bait.Em.Emit("command.success", nil)
		}
	}
}

// Run command in sh interpreter
func execCommand(cmd string) error {
	file, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}
	runner, _ := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
	)
	err = runner.Run(context.TODO(), file)

	return err
}

func splitInput(input string) ([]string, string) {
	// end my suffering
	// TODO: refactor this garbage
	quoted := false
	startlastcmd := false
	lastcmddone := false
	cmdArgs := []string{}
	sb := &strings.Builder{}
	cmdstr := &strings.Builder{}
	lastcmd := readline.GetHistory(readline.HistorySize() - 1)

	for _, r := range input {
		if r == '"' {
			// start quoted input
			// this determines if other runes are replaced
			quoted = !quoted
			// dont add back quotes
			//sb.WriteRune(r)
		} else if !quoted && r == '~' {
			// if not in quotes and ~ is found then make it $HOME
			sb.WriteString(os.Getenv("HOME"))
		} else if !quoted && r == ' ' {
			// if not quoted and there's a space then add to cmdargs
			cmdArgs = append(cmdArgs, sb.String())
			sb.Reset()
		} else if !quoted && r == '^' && startlastcmd && !lastcmddone {
			// if ^ is found, isnt in quotes and is
			// the second occurence of the character and is
			// the first time "^^" has been used
			cmdstr.WriteString(lastcmd)
			sb.WriteString(lastcmd)

			startlastcmd = !startlastcmd
			lastcmddone = !lastcmddone

			continue
		} else if !quoted && r == '^' && !lastcmddone {
			// if ^ is found, isnt in quotes and is the
			// first time of starting "^^"
			startlastcmd = !startlastcmd
			continue
		} else {
			sb.WriteRune(r)
		}
		cmdstr.WriteRune(r)
	}
	if sb.Len() > 0 {
		cmdArgs = append(cmdArgs, sb.String())
	}

	readline.AddHistory(input)
	readline.SaveHistory(homedir + "/.hilbish-history")
	return cmdArgs, cmdstr.String()
}

func StartMultiline(prev string, sb *strings.Builder) bool {
	// sb fromt outside is passed so we can
	// save input from previous prompts
	if sb.String() == "" { sb.WriteString(prev + "\n") }

	fmt.Printf("... ")

	reader := bufio.NewReader(os.Stdin)

	cont, err := reader.ReadString('\n')
	if err == io.EOF {
		// Exit when ^D
		fmt.Println("")
		return true
	}

	sb.WriteString(cont)

	err = execCommand(sb.String())
	if err != nil && syntax.IsIncomplete(err) {
		return false
	}

	return true
}
