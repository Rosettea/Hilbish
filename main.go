package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"os/signal"
	"strings"
)

func main() {
	HandleSignals()

	for {
		user, _ := user.Current()
		dir, _ := os.Getwd()
		host, _ := os.Hostname()

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("\u001b[1m\u001b[36m%s@%s \u001b[34m%s $ \u001b[0m", user.Username, host, dir)

		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		cmdString = strings.TrimSuffix(cmdString, "\n")
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
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
	}()
}
