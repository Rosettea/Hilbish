package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func RunInput(input string) {
	// First try to load input, essentially compiling to bytecode
	fn, err := l.LoadString(input)
	if err != nil && noexecute {
		fmt.Println(err)
		return
	}
	// And if there's no syntax errors and -n isnt provided, run
	if !noexecute {
		l.Push(fn)
		err = l.PCall(0, lua.MultRet, nil)
	}
	if err == nil {
		// If it succeeds, add to history and prompt again
		HandleHistory(input)

		hooks.Em.Emit("command.exit", 0)
		return
	}

	cmdArgs, cmdString := splitInput(input)

	// If alias was found, use command alias
	if aliases[cmdArgs[0]] != "" {
		alias := aliases[cmdArgs[0]]
		cmdString = alias + strings.Trim(cmdString, cmdArgs[0])
		cmdArgs, cmdString = splitInput(cmdString)
	}

	// If command is defined in Lua then run it
	if commands[cmdArgs[0]] {
		err := l.CallByParam(lua.P{
			Fn: l.GetField(
				l.GetTable(
					l.GetGlobal("commanding"),
					lua.LString("__commands")),
				cmdArgs[0]),
			NRet:    1,
			Protect: true,
		}, luar.New(l, cmdArgs[1:]))
		luaexitcode := l.Get(-1)
		exitcode := lua.LNumber(0)

		l.Pop(1)

		if code, ok := luaexitcode.(lua.LNumber); luaexitcode != lua.LNil && ok {
			exitcode = code
		}
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Error in command:\n\n" + err.Error())
		}
		if cmdArgs[0] != "exit" {
			HandleHistory(cmdString)
		}
		hooks.Em.Emit("command.exit", exitcode)
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
					hooks.Em.Emit("command.exit", code)
				} else if err != nil {
					fmt.Fprintln(os.Stderr, err)
					hooks.Em.Emit("command.exit", 1)
				}
				break
			}
		} else {
			if code, ok := interp.IsExitStatus(err); ok {
				hooks.Em.Emit("command.exit", code)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	} else {
		hooks.Em.Emit("command.exit", 0)
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

	return cmdArgs, cmdstr.String()
}

func HandleHistory(cmd string) {
	readline.AddHistory(cmd)
	readline.SaveHistory(homedir + "/.hilbish-history")
	// TODO: load history again (history shared between sessions like this ye)
}

