package main

import (
	_ "bufio"
	"fmt"
	"os"
	"os/exec"
	_ "os/user"
	"syscall"
	"os/signal"
	"strings"
	"io"
	lfs "hilbish/golibs/fs"
	cmds "hilbish/golibs/commander"

	"github.com/akamensky/argparse"
	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

const version = "0.0.6"
var l *lua.LState
var prompt string
var commands = map[string]bool{}

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
	HandleSignals()
	LuaInit()

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

		if err == nil { continue }

		quoted := false
		q := func(r rune) bool {
			if r == '"' {
				quoted = !quoted
			}
			return !quoted && r == ' '
		}

		cmdArgs := strings.FieldsFunc(cmdString, q)

		if len(cmdArgs) == 0 { continue }

		if commands[cmdArgs[0]] {
			err := l.CallByParam(lua.P{
				Fn: l.GetField(
					l.GetTable(
						l.GetGlobal("commander"),
						lua.LString("__commands")),
					cmdArgs[0]),
				NRet: 0,
				Protect: true,
			}, luar.New(l, cmdArgs[1:]))
			if err != nil {
				// TODO: dont panic
				panic(err)
			}
			continue
		}
		switch cmdArgs[0] {
		case "exit":
			os.Exit(0)
		default:
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		readline.AddHistory(cmdString)
	}
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

	l.PreloadModule("fs", lfs.Loader)

	commander := cmds.New()
	commander.Events.On("commandRegister",
	func (cmdName string, cmd *lua.LFunction) {
		commands[cmdName] = true
		l.SetField(
			l.GetTable(l.GetGlobal("commander"),
			lua.LString("__commands")),
			cmdName,
			cmd)
	})

	l.PreloadModule("commander", commander.Loader)

	err := l.DoFile(os.Getenv("HOME") + "/.hilbishrc.lua")
	if err != nil {
		panic(err)
	}
}
