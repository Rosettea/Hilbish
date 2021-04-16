package main

import (
	"fmt"
	"os"
	"os/user"
	"os/signal"
	"strings"
	"io"
	hooks "hilbish/golibs/bait"

	"github.com/akamensky/argparse"
	"github.com/bobappleyard/readline"
	"github.com/Hilbis/Hilbiline"
	"github.com/yuin/gopher-lua"
	"golang.org/x/term"
)

const version = "0.3.3-hilbiline"
var l *lua.LState
// User's prompt, this will get set when lua side is initialized
var prompt string
var multilinePrompt = "> "

// Map of builtin/custom commands defined in the commander lua module
var commands = map[string]bool{}
// Command aliases
var aliases = map[string]string{}
var bait hooks.Bait
var homedir string
var running bool

func main() {
	homedir, _ = os.UserHomeDir()
	defaultconfpath := homedir + "/.hilbishrc.lua"

	parser := argparse.NewParser("hilbish", "A shell for lua and flower lovers")
	verflag := parser.Flag("v", "version", &argparse.Options{
		Required: false,
		Help: "Prints Hilbish version",
	})
	setshflag := parser.Flag("S", "set-shell-env", &argparse.Options{
		Required: false,
		Help: "Sets $SHELL to Hilbish's executed path",
	})
	configflag := parser.String("C", "config", &argparse.Options{
		Required: false,
		Help: "Sets the path to Hilbish's config",
		Default: defaultconfpath,
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

	// If user's config doesn't exixt,
	if _, err := os.Stat(defaultconfpath); os.IsNotExist(err) {
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

	go HandleSignals()
	LuaInit(*configflag)

	hl := hilbiline.New(prompt)
	//readline.Completer = readline.FilenameCompleter
	//readline.LoadHistory(homedir + "/.hilbish-history")

	for {
		running = false

		hl.SetPrompt(fmtPrompt())
		input, err := hl.Read()

		if err == io.EOF {
			// Exit if user presses ^D (ctrl + d)
			fmt.Println("")
			break
		}
		if err != nil {
			// If we get a completely random error, print
			fmt.Fprintln(os.Stderr, err)
		}

		input = strings.TrimSpace(input)
		if len(input) == 0 { fmt.Print("\n"); continue }

		if strings.HasSuffix(input, "\\") {
			for {
				input, err = ContinuePrompt(strings.TrimSuffix(input, "\\"))

				if err != nil || !strings.HasSuffix(input, "\\") { break }
			}
		}
		running = true
		RunInput(input)

		termwidth, _, err := term.GetSize(0)
		if err != nil { continue }
		fmt.Printf("\u001b[7m∆\u001b[0m" + strings.Repeat(" ", termwidth - 1) + "\r")
	}
}

func ContinuePrompt(prev string) (string, error) {
	cont, err := readline.String(multilinePrompt)
	if err != nil {
		fmt.Println("")
		return "", err
	}
	cont = strings.TrimSpace(cont)

	return prev + "\n" + strings.TrimSuffix(cont, "\n"), nil
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
		"u", user.Username,
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
	signal.Notify(c, os.Interrupt)

	for range c {
		if !running {
			//fmt.Println(" // interrupt")
			//readline.ReplaceLine("", 0)
			//readline.RefreshLine()
		}
	}
}

