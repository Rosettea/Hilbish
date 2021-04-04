package main

import (
	"fmt"
	"os"
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/bobappleyard/readline"
	_ "github.com/yuin/gopher-lua"
	_ "layeh.com/gopher-luar"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

)

func RunInput(input string) {
	// First try to run user input in Lua
	if strings.HasSuffix(input, "\\") {
		for {
			input = strings.TrimSuffix(input, "\\")
			input = ContinuePrompt(input)
			input = strings.TrimSpace(input)
			if input == "" { break }
			// For some reason !HasSuffix didnt work :\, stupid
			if !strings.HasSuffix(input, "\\") { break }
		}
	}

	err := l.DoString(input)

	if err == nil {
		// If it succeeds, add to history and prompt again
		//readline.AddHistory(input)
		//readline.SaveHistory(homedir + "/.hilbish-history")
		bait.Em.Emit("command.exit", nil)
		bait.Em.Emit("command.success", nil)
		return
	}
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

	readline.AddHistory(input)
	readline.SaveHistory(homedir + "/.hilbish-history")
	return cmdArgs, cmdstr.String()
}

func ContinuePrompt(prev string) string {
	fmt.Printf("... ")

	reader := bufio.NewReader(os.Stdin)

	cont, err := reader.ReadString('\n')
	if err == io.EOF {
		// Exit when ^D
		fmt.Println("")
		return ""
	}

	return prev + "\n" + cont
}

func StartMultiline(prev string, sb *strings.Builder) bool {
	// sb fromt outside is passed so we can
	// save input from previous prompts
	if sb.String() == "" { sb.WriteString(prev + "\n") }

	fmt.Printf("... ")

	reader := bufio.NewReader(os.Stdin)

	cont, err := reader.ReadString('\n')
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
