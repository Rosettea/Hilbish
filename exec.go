package main

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"hilbish/util"

	"github.com/yuin/gopher-lua"
	"mvdan.cc/sh/v3/shell"
	//"github.com/yuin/gopher-lua/parse"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"mvdan.cc/sh/v3/expand"
)

var errNotExec = errors.New("not executable")
var runnerMode lua.LValue = lua.LString("hybrid")

func runInput(input string, priv bool) {
	running = true
	cmdString := aliases.Resolve(input)
	hooks.Em.Emit("command.preexec", input, cmdString)

	if runnerMode.Type() == lua.LTString {
		switch runnerMode.String() {
			case "hybrid":
				_, err := handleLua(cmdString)
				if err == nil {
					cmdFinish(0, cmdString, priv)
					return
				}
				exitCode, err := handleSh(cmdString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				cmdFinish(exitCode, cmdString, priv)
			case "hybridRev":
				_, err := handleSh(cmdString)
				if err == nil {
					cmdFinish(0, cmdString, priv)
					return
				}
				exitCode, err := handleLua(cmdString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				cmdFinish(exitCode, cmdString, priv)
			case "lua":
				exitCode, err := handleLua(cmdString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				cmdFinish(exitCode, cmdString, priv)
			case "sh":
				exitCode, err := handleSh(cmdString)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				cmdFinish(exitCode, cmdString, priv)
		}
	} else {
		// can only be a string or function so
		err := l.CallByParam(lua.P{
			Fn: runnerMode,
			NRet: 2,
			Protect: true,
		}, lua.LString(cmdString))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			cmdFinish(124, cmdString, priv)
			return
		}

		luaexitcode := l.Get(-2) // first return value (makes sense right i love stacks)
		runErr := l.Get(-1)
		l.Pop(2)

		var exitCode uint8
		if code, ok := luaexitcode.(lua.LNumber); luaexitcode != lua.LNil && ok {
			exitCode = uint8(code)
		}

		if runErr != lua.LNil {
			fmt.Fprintln(os.Stderr, runErr)
		}
		cmdFinish(exitCode, cmdString, priv)
	}
}

func handleLua(cmdString string) (uint8, error) {
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
		return 125, err
	}
	// And if there's no syntax errors and -n isnt provided, run
	if !noexecute {
		l.Push(fn)
		err = l.PCall(0, lua.MultRet, nil)
	}
	if err == nil {
		return 0, nil
	}

	return 125, err
}

func handleSh(cmdString string) (uint8, error) {
	err := execCommand(cmdString)
	if err != nil {
		// If input is incomplete, start multiline prompting
		if syntax.IsIncomplete(err) {
			for {
				cmdString, err = continuePrompt(strings.TrimSuffix(cmdString, "\\"))
				if err != nil {
					break
				}
				err = execCommand(cmdString)
				if syntax.IsIncomplete(err) || strings.HasSuffix(cmdString, "\\") {
					continue
				} else if code, ok := interp.IsExitStatus(err); ok {
					return code, nil
				} else if err != nil {
					return 126, err
				} else {
					return 0, nil
				}
			}
		} else {
			if code, ok := interp.IsExitStatus(err); ok {
				return code, nil
			} else {
				return 126, err
			}
		}
	}

	return 0, nil
}

// Run command in sh interpreter
func execCommand(cmd string) error {
	file, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}

	var bg bool
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

			return interp.NewExitStatus(exitcode)
		}

		err := lookpath(args[0])
		if err == errNotExec {
			hooks.Em.Emit("command.no-perm", args[0])
			hooks.Em.Emit("command.not-executable", args[0])
			return interp.NewExitStatus(126)
		} else if err != nil {
			hooks.Em.Emit("command.not-found", args[0])
			return interp.NewExitStatus(127)
		}

		killTimeout := 2 * time.Second
		// from here is basically copy-paste of the default exec handler from
		// sh/interp but with our job handling
		hc := interp.HandlerCtx(ctx)
		path, err := interp.LookPathDir(hc.Dir, hc.Env, args[0])
		if err != nil {
			fmt.Fprintln(hc.Stderr, err)
			return interp.NewExitStatus(127)
		}

		env := hc.Env
		envList := make([]string, 0, 64)
		env.Each(func(name string, vr expand.Variable) bool {
			if !vr.IsSet() {
				// If a variable is set globally but unset in the
				// runner, we need to ensure it's not part of the final
				// list. Seems like zeroing the element is enough.
				// This is a linear search, but this scenario should be
				// rare, and the number of variables shouldn't be large.
				for i, kv := range envList {
					if strings.HasPrefix(kv, name+"=") {
						envList[i] = ""
					}
				}
			}
			if vr.Exported && vr.Kind == expand.String {
				envList = append(envList, name+"="+vr.String())
			}
			return true
		})
		cmd := exec.Cmd{
			Path: path,
			Args: args,
			Env: envList,
			Dir: hc.Dir,
			Stdin: hc.Stdin,
			Stdout: hc.Stdout,
			Stderr: hc.Stderr,
		}

		err = cmd.Start()
		job := jobs.getLatest()
		if err == nil {
			if bg {
				job.start(cmd.Process.Pid)
			}

			if done := ctx.Done(); done != nil {
				go func() {
					<-done

					if killTimeout <= 0 || runtime.GOOS == "windows" {
						cmd.Process.Signal(os.Kill)
						return
					}

					// TODO: don't temporarily leak this goroutine
					// if the program stops itself with the
					// interrupt.
					go func() {
						time.Sleep(killTimeout)
						cmd.Process.Signal(os.Kill)
					}()
					cmd.Process.Signal(os.Interrupt)
				}()
			}

			err = cmd.Wait()
		}

		var exit uint8
		switch x := err.(type) {
		case *exec.ExitError:
			// started, but errored - default to 1 if OS
			// doesn't have exit statuses
			if status, ok := x.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					if ctx.Err() != nil {
						return ctx.Err()
					}
					exit = uint8(128 + status.Signal())
					goto end
				}
				exit = uint8(status.ExitStatus())
				goto end
			}
			exit = 1
			goto end
		case *exec.Error:
			// did not start
			fmt.Fprintf(hc.Stderr, "%v\n", err)
			exit = 127
			goto end
		case nil:
			goto end
		default:
			return err
		}
		end:
		if bg {
			job.exitCode = int(exit)
			job.finish()
		}
		return interp.NewExitStatus(exit)
	}

	runner, _ := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.ExecHandler(exechandle),
	)

	buf := new(bytes.Buffer)
	printer := syntax.NewPrinter()

	for _, stmt := range file.Stmts {
		bg = false
		if stmt.Background {
			bg = true
			printer.Print(buf, stmt.Cmd)

			stmtStr := buf.String()
			buf.Reset()
			jobs.add(stmtStr)
		}

		err = runner.Run(context.TODO(), stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func lookpath(file string) error { // custom lookpath function so we know if a command is found *and* is executable
	var skip []string
	if runtime.GOOS == "windows" {
		skip = []string{"./", "../", "~/", "C:"}
	} else {
		skip = []string{"./", "/", "../", "~/"}
	}
	for _, s := range skip {
		if strings.HasPrefix(file, s) {
			return findExecutable(file, false, false)
		}
	}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		path := filepath.Join(dir, file)
		err := findExecutable(path, true, false)
		if err == errNotExec {
			return err
		} else if err == nil {
			return nil
		}
	}

	return os.ErrNotExist
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

func cmdFinish(code uint8, cmdstr string, private bool) {
	// if input has space at the beginning, dont put in history
	if interactive && !private {
		handleHistory(cmdstr)
	}
	util.SetField(l, hshMod, "exitCode", lua.LNumber(code), "Exit code of last exected command")
	hooks.Em.Emit("command.exit", code, cmdstr)
}
