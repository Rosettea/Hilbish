package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/gopher-lua"
	"mvdan.cc/sh/v3/shell"
	//"github.com/yuin/gopher-lua/parse"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func runInput(input string) {
	running = true
	cmdString := aliases.Resolve(input)

	hooks.Em.Emit("command.preexec", input, cmdString)

	// First try to load input, essentially compiling to bytecode
	fn, err := l.LoadString(cmdString)
	if err != nil && noexecute {
		fmt.Println(err)
	/*	if lerr, ok := err.(*lua.ApiError); ok {
			if perr, ok := lerr.Cause.(*parse.Error); ok {
				print(perr.Pos.Line == parse.EOF)
			}
		}
	*/
		return
	}
	// And if there's no syntax errors and -n isnt provided, run
	if !noexecute {
		l.Push(fn)
		err = l.PCall(0, lua.MultRet, nil)
	}
	if err == nil {
		cmdFinish(0, cmdString)
		return
	}

	// Last option: use sh interpreter
	err = execCommand(cmdString)
	if err != nil {
		// If input is incomplete, start multiline prompting
		if syntax.IsIncomplete(err) {
			for {
				cmdString, err = continuePrompt(strings.TrimSuffix(cmdString, "\\"))
				if err != nil {
					break
				}
				err = execCommand(cmdString)
				if syntax.IsIncomplete(err) || strings.HasSuffix(input, "\\") {
					continue
				} else if code, ok := interp.IsExitStatus(err); ok {
					cmdFinish(code, cmdString)
				} else if err != nil {
					fmt.Fprintln(os.Stderr, err)
					cmdFinish(1, cmdString)
				}
				break
			}
		} else {
			if code, ok := interp.IsExitStatus(err); ok {
				cmdFinish(code, cmdString)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	} else {
		cmdFinish(0, cmdString)
	}
}

// Run command in sh interpreter
func execCommand(cmd string) error {
	file, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}

	exechandle := func(ctx context.Context, args []string) error {
		_, argstring := splitInput(strings.Join(args, " "))
		// i dont really like this but it works
		if aliases.All()[args[0]] != "" {
			for i, arg := range args {
				if strings.Contains(arg, " ") {
					args[i] = fmt.Sprintf("\"%s\"", arg)
				}
			}
			_, argstring = splitInput(strings.Join(args, " "))

			// If alias was found, use command alias
			argstring = aliases.Resolve(argstring)
			args, _ = shell.Fields(argstring, nil)
		}

		// If command is defined in Lua then run it
		luacmdArgs := l.NewTable()
		for _, str := range args[1:] {
			luacmdArgs.Append(lua.LString(str))
		}

		if commands[args[0]] != nil {
			err := l.CallByParam(lua.P{
				Fn: commands[args[0]],
				NRet:    1,
				Protect: true,
			}, luacmdArgs)
			
			if err != nil {
				fmt.Fprintln(os.Stderr,
					"Error in command:\n\n" + err.Error())
				return interp.NewExitStatus(1)
			}

			luaexitcode := l.Get(-1)
			var exitcode uint8

			l.Pop(1)

			if code, ok := luaexitcode.(lua.LNumber); luaexitcode != lua.LNil && ok {
				exitcode = uint8(code)
			}

			cmdFinish(exitcode, argstring)
			return interp.NewExitStatus(exitcode)
		}

		err := lookpath(args[0])
		if err == os.ErrPermission {
			hooks.Em.Emit("command.no-perm", args[0])
			return interp.NewExitStatus(126)
		} else if err != nil {
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

// custom lookpath function so we know if a command is found *and* has execute permission
func lookpath(file string) error {
	skip := []string{"./", "/", "../", "~/"}
	for _, s := range skip {
		if strings.HasPrefix(file, s) {
			err := findExecutable(file)
			return err
		}
	}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		path := filepath.Join(dir, file)
		err := findExecutable(path)
		if err == os.ErrPermission {
			return err
		} else if err == nil {
			return nil
		}
	}

	return os.ErrNotExist
}

func findExecutable(name string) error {
	f, err := os.Stat(name)
	if err != nil {
		return err
	}
	if m := f.Mode(); !m.IsDir() && m & 0111 != 0 {
		return nil
	}
	return os.ErrPermission
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

func cmdFinish(code uint8, cmdstr string) {
	hooks.Em.Emit("command.exit", code, cmdstr)
}
