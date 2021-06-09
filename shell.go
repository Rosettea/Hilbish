package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

//	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"layeh.com/gopher-luar"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func RunInput(input string) {
	cmdArgs, cmdString := splitInput(input)

	// If alias was found, use command alias
	for aliases[cmdArgs[0]] != "" {
		alias := aliases[cmdArgs[0]]
		cmdString = alias + strings.TrimPrefix(cmdString, cmdArgs[0])
		cmdArgs, cmdString = splitInput(cmdString)

		if aliases[cmdArgs[0]] != "" {
			continue
		}
	}

	// First try to load input, essentially compiling to bytecode
	fn, err := l.LoadString(cmdString)
	if err != nil && noexecute {
		fmt.Println(err)
		if lerr, ok := err.(*lua.ApiError); ok {
			if perr, ok := lerr.Cause.(*parse.Error); ok {
				print(perr.Pos.Line == parse.EOF)
			}
		}
		return
	}
	// And if there's no syntax errors and -n isnt provided, run
	if !noexecute {
		l.Push(fn)
		err = l.PCall(0, lua.MultRet, nil)
	}
	if err == nil {
		hooks.Em.Emit("command.exit", 0)
		return
	}
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
		var exitcode uint8 = 0

		l.Pop(1)

		if code, ok := luaexitcode.(lua.LNumber); luaexitcode != lua.LNil && ok {
			exitcode = uint8(code)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Error in command:\n\n" + err.Error())
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
}

// Run command in sh interpreter
func execCommand(cmd string) error {
	file, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}

	exechandle := func(ctx context.Context, args []string) error {
		hc := interp.HandlerCtx(ctx)
		_, argstring := splitInput(strings.Join(args, " "))

		// If alias was found, use command alias
		if aliases[args[0]] != "" {
			alias := aliases[args[0]]
			argstring = alias + strings.TrimPrefix(argstring, args[0])
			args[0] = alias
		}

		// If command is defined in Lua then run it
		if commands[args[0]] {
			err := l.CallByParam(lua.P{
				Fn: l.GetField(
					l.GetTable(
						l.GetGlobal("commanding"),
						lua.LString("__commands")),
					args[0]),
				NRet:    1,
				Protect: true,
			}, luar.New(l, args[1:]))
			luaexitcode := l.Get(-1)
			var exitcode uint8 = 0

			l.Pop(1)

			if code, ok := luaexitcode.(lua.LNumber); luaexitcode != lua.LNil && ok {
				exitcode = uint8(code)
			}

			if err != nil {
				fmt.Fprintln(os.Stderr,
					"Error in command:\n\n" + err.Error())
			}
			hooks.Em.Emit("command.exit", exitcode)
			return interp.NewExitStatus(exitcode)
		}

		if _, err := interp.LookPathDir(hc.Dir, hc.Env, args[0]); err != nil {
			hooks.Em.Emit("command.not-found", args[0])
			return interp.NewExitStatus(127)
		}

		return interp.DefaultExecHandler(2 * time.Second)(ctx, args)
	}
	runner, _ := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.ExecHandler(exechandle),
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

