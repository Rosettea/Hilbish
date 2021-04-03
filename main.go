package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
	"os/signal"
	"strings"
	"io"
	hooks "hilbish/golibs/bait"

	"github.com/akamensky/argparse"
	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"

)

const version = "0.3.0-dev"
var l *lua.LState
// User's prompt, this will get set when lua side is initialized
var prompt string
// Map of builtin/custom commands defined in the commander lua module
var commands = map[string]bool{}
// Command aliases
var aliases = map[string]string{}
var bait hooks.Bait
var homedir string

func main() {
	parser := argparse.NewParser("hilbish", "A shell for lua and flower lovers")
	verflag := parser.Flag("v", "version", &argparse.Options{
		Required: false,
		Help: "Prints Hilbish version",
	})
	setshflag := parser.Flag("S", "set-shell-env", &argparse.Options{
		Required: false,
		Help: "Sets $SHELL to Hilbish's executed path",
	})

	err := parser.Parse(os.Args)
	// If invalid flags or --help/-h,
	if err != nil {
		// Print usage
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	if *verflag {
		fmt.Printf("Hilbish v%s\n", version)
		os.Exit(0)
	}

	// Set $SHELL if the user wants to
	if *setshflag { os.Setenv("SHELL", os.Args[0]) }

	homedir, _ = os.UserHomeDir()
	// If user's config doesn't exixt,
	if _, err := os.Stat(homedir + "/.hilbishrc.lua"); os.IsNotExist(err) {
		// Read default from current directory
		// (this is assuming the current dir is Hilbish's git)
		input, err := os.ReadFile(".hilbishrc.lua")
		if err != nil {
			// If it wasnt found, go to "real default"
			input, err = os.ReadFile("/usr/share/hilbish/.hilbishrc.lua")
			if err != nil {
				fmt.Println("could not find .hilbishrc.lua or /usr/share/hilbish/.hilbishrc.lua")
				return
			}
		}

		// Create it using either default config we found
		err = os.WriteFile(homedir + "/.hilbishrc.lua", input, 0644)
		if err != nil {
			// If that fails, bail
			fmt.Println("Error creating config file")
			fmt.Println(err)
			return
		}
	}

	HandleSignals()
	LuaInit()

	readline.Completer = readline.FilenameCompleter
	readline.LoadHistory(homedir + "/.hilbish-history")

	for {
		cmdString, err := readline.String(fmtPrompt())
		if err == io.EOF {
			// Exit if user presses ^D (ctrl + d)
			fmt.Println("")
			break
		}
		if err != nil {
			// If we get a completely random error, print
			fmt.Fprintln(os.Stderr, err)
		}

		// I have no idea if we need this anymore
		cmdString = strings.TrimSuffix(cmdString, "\n")
		RunInput(cmdString)
	}
}

// This semi cursed function formats our prompt (obviously)
func fmtPrompt() string {
	user, _ := user.Current()
	host, _ := os.Hostname()
	cwd, _ := os.Getwd()

	cwd = strings.Replace(cwd, user.HomeDir, "~", 1)

	args := []string{
		"d", cwd,
		"h", host,
		"u", user.Name,
	}

	for i, v := range args {
		if i % 2 == 0 {
			args[i] = "%" + v
		}
	}

	r := strings.NewReplacer(args...)
	nprompt := r.Replace(prompt)

	return nprompt
}

// do i even have to say
func HandleSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()
}

