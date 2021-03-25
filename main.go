package main

import (
	"bufio"
	"fmt"
	"os"
	_ "os/exec"
	_ "os/user"
	"syscall"
	"os/signal"
	"strings"
	"io"
	"context"
	lfs "hilbish/golibs/fs"
	cmds "hilbish/golibs/commander"

	"github.com/akamensky/argparse"
	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

)

const version = "0.1.2"
var l *lua.LState
var prompt string
var commands = map[string]bool{}
var aliases = map[string]string{}

func main() {
	parser := argparse.NewParser("hilbish", "A shell for lua and flower lovers")
	verflag := parser.Flag("v", "version", &argparse.Options{
		Required: false,
		Help: "prints hilbish version",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	if *verflag {
		fmt.Printf("Hilbish v%s\n", version)
		os.Exit(0)
	}

	os.Setenv("SHELL", os.Args[0])

	input, err := os.ReadFile(".hilbishrc.lua")
	if err != nil {
		input, err = os.ReadFile("/usr/share/hilbish/.hilbishrc.lua")
		if err != nil {
			fmt.Println("could not find .hilbishrc.lua or /usr/share/hilbish/.hilbishrc.lua")
			return
		}
	}

	homedir, _ := os.UserHomeDir()
	if _, err := os.Stat(homedir + "/.hilbishrc.lua"); os.IsNotExist(err) {
		err = os.WriteFile(homedir + "/.hilbishrc.lua", input, 0644)
		if err != nil {
			fmt.Println("Error creating config file")
			fmt.Println(err)
			return
		}
	}

	HandleSignals()
	LuaInit()

	readline.Completer = readline.FilenameCompleter
	for {
		//user, _ := user.Current()
		//dir, _ := os.Getwd()
		//host, _ := os.Hostname()

		//reader := bufio.NewReader(os.Stdin)

		//fmt.Printf(prompt)

		cmdString, err := readline.String(prompt)
		if err == io.EOF {
			fmt.Println("")
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		cmdString = strings.TrimSuffix(cmdString, "\n")
		err = l.DoString(cmdString)

		if err == nil {
			readline.AddHistory(cmdString)
			continue
		}

		cmdArgs := splitInput(cmdString)
		if len(cmdArgs) == 0 { continue }

		if aliases[cmdArgs[0]] != "" {
			cmdString = aliases[cmdArgs[0]] + strings.Trim(cmdString, cmdArgs[0])
			//cmdArgs := splitInput(cmdString)
			execCommand(cmdString)
			continue
		}

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
			continue
		}
		switch cmdArgs[0] {
		case "exit":
			os.Exit(0)
		default:
			err := execCommand(cmdString)
			if err != nil {
				if syntax.IsIncomplete(err) {
					sb := &strings.Builder{}
					for {
						done := StartMultiline(cmdString, sb)
						if done {
							break
						}
					}
				} else {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}

func StartMultiline(prev string, sb *strings.Builder) bool {
	if sb.String() == "" { sb.WriteString(prev + "\n") }

	fmt.Printf("... ")

	reader := bufio.NewReader(os.Stdin)

	cont, err := reader.ReadString('\n')
	if err == io.EOF {
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

func splitInput(input string) []string {
	quoted := false
	cmdArgs := []string{}
	sb := &strings.Builder{}

	for _, r := range input {
		if r == '"' {
			quoted = !quoted
			// dont add back quotes
			//sb.WriteRune(r)
		} else if !quoted && r == '~' {
			sb.WriteString(os.Getenv("HOME"))
		} else if !quoted && r == ' ' {
			cmdArgs = append(cmdArgs, sb.String())
			sb.Reset()
		} else {
			sb.WriteRune(r)
		}
	}
	if sb.Len() > 0 {
		cmdArgs = append(cmdArgs, sb.String())
	}

	readline.AddHistory(input)
	return cmdArgs
}

func execCommand(cmd string) error {
	file, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return err
	}
	runner, _ := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
	)
	runner.Run(context.TODO(), file)

	return nil
}

func HandleSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()
}

func LuaInit() {
	l = lua.NewState()

	l.OpenLibs()

	l.SetGlobal("prompt", l.NewFunction(hshprompt))
	l.SetGlobal("alias", l.NewFunction(hshalias))

	l.PreloadModule("fs", lfs.Loader)

	commander := cmds.New()
	commander.Events.On("commandRegister",
	func (cmdName string, cmd *lua.LFunction) {
		commands[cmdName] = true
		l.SetField(
			l.GetTable(l.GetGlobal("commanding"),
			lua.LString("__commands")),
			cmdName,
			cmd)
	})

	l.PreloadModule("commander", commander.Loader)

	l.DoString("package.path = package.path .. ';./libs/?/init.lua;/usr/share/hilbish/libs/?/init.lua'")

	err := l.DoFile("/usr/share/hilbish/preload.lua")
	if err != nil {
		err = l.DoFile("preload.lua")
		if err != nil {
			fmt.Fprintln(os.Stderr,
			"Missing preload file, builtins may be missing.")
		}
	}

	homedir, _ := os.UserHomeDir()
	err = l.DoFile(homedir + "/.hilbishrc.lua")
	if err != nil {
		panic(err)
	}
}
