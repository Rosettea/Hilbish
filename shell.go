package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

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
		HandleHistory(input)

		bait.Em.Emit("command.exit", 0)
		return
	}

	cmdArgs, cmdString := splitInput(input)

	// If alias was found, use command alias
	if aliases[cmdArgs[0]] != "" {
		alias := aliases[cmdArgs[0]]
		cmdString = alias + strings.Trim(cmdString, cmdArgs[0])
		cmdArgs[0] = alias
	}

	// If command is defined in Lua then run it
	if commands[cmdArgs[0]] {
		err := l.CallByParam(lua.P{
			Fn: l.GetField(
				l.GetTable(
					l.GetGlobal("commanding"),
					lua.LString("__commands")),
				cmdArgs[0]),
			NRet:    0,
			Protect: true,
		}, luar.New(l, cmdArgs[1:]))
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Error in command:\n\n"+err.Error())
		}
		if cmdArgs[0] != "exit" {
			HandleHistory(cmdString)
		}
		return
	}

	// Last option: use sh interpreter
	err = execCommand(cmdString)
	if err != nil {
		// If input is incomplete, start multiline prompting
		if syntax.IsIncomplete(err) {
			for {
				cmdString, err = ContinuePrompt(strings.TrimSuffix(cmdString, "\\"))
				if err != nil {
					break
				}
				err = execCommand(cmdString)
				if syntax.IsIncomplete(err) || strings.HasSuffix(input, "\\") {
					continue
				} else if code, ok := interp.IsExitStatus(err); ok {
					bait.Em.Emit("command.exit", code)
				} else if err != nil {
					fmt.Fprintln(os.Stderr, err)
					bait.Em.Emit("command.exit", 1)
				}
				break
			}
		} else {
			if code, ok := interp.IsExitStatus(err); ok {
				bait.Em.Emit("command.exit", code)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	} else {
		bait.Em.Emit("command.exit", 0)
	}
	HandleHistory(cmdString)
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
	lastcmd := "" //readline.GetHistory(readline.HistorySize() - 1)

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

	return cmdArgs, cmdstr.String()
}

func HandleHistory(cmd string) {
	//readline.AddHistory(cmd)
	//readline.SaveHistory(homedir + "/.hilbish-history")
	// TODO: load history again (history shared between sessions like this ye)
}

func StartMultiline(prev string, sb *strings.Builder) bool {
	// sb fromt outside is passed so we can
	// save input from previous prompts
	if sb.String() == "" { sb.WriteString(prev + " ") }

	fmt.Print(multilinePrompt)

	reader := bufio.NewReader(os.Stdin)

	cont, err := reader.ReadString('\n')
	cont = strings.TrimSuffix(cont, "\n") + " "
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
