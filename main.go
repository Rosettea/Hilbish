package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	_ "os/user"
	"os/signal"
	"strings"
	"github.com/yuin/gopher-lua"
)

var l *lua.LState
var prompt string

func main() {
	HandleSignals()
	LuaInit()

	for {
		//user, _ := user.Current()
		//dir, _ := os.Getwd()
		//host, _ := os.Hostname()

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf(prompt)

		cmdString, err := reader.ReadString('\n')
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
		}

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func HandleSignals() {
	signal.Ignore(os.Interrupt)
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
