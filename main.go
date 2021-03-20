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

	"github.com/akamensky/argparse"
	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
)

const version = "0.0.4"
var l *lua.LState
var prompt string

func main() {
	parser := argparse.NewParser("hilbish", "A shell for lua and flower lovers")
	verflag := parser.Flag("v", "version", &argparse.Options{
		Required: false,
		Help: "color palette to use",
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
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		cmdString = strings.TrimSuffix(cmdString, "\n")
		err = l.DoString(cmdString)

		if err == nil { continue }

		cmdArgs := strings.Fields(cmdString)

		if len(cmdArgs) == 0 { continue }

		switch cmdArgs[0] {
		case "exit":
			os.Exit(0)
		case "cd":
			os.Chdir(strings.Trim(cmdString, "cd "))
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

	err := l.DoFile(os.Getenv("HOME") + "/.hilbishrc.lua")
	if err != nil {
		panic(err)
	}
}
