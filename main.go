package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"os/signal"
	"os/user"
	"path/filepath"

	"hilbish/golibs/bait"

	"github.com/pborman/getopt"
	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
	"golang.org/x/term"
)

const version = "0.4.0"

var (
	l *lua.LState

	prompt string // User's prompt, this will get set when lua side is initialized
	multilinePrompt = "> "

	commands = map[string]bool{}
	aliases = map[string]string{}

	homedir string
	curuser *user.User

	running bool // Is a command currently running
	hooks bait.Bait
	interactive bool
	login bool // Are we the login shell?
	noexecute bool // Should we run Lua or only report syntax errors
)

func main() {
	homedir, _ = os.UserHomeDir()
	curuser, _ = user.Current()
	defaultconfpath := homedir + "/.hilbishrc.lua"

//	parser := argparse.NewParser("hilbish", "A shell for lua and flower lovers")
	verflag := getopt.BoolLong("version", 'v', "Prints Hilbish version")
	setshflag := getopt.BoolLong("setshellenv", 'S', "Sets $SHELL to Hilbish's executed path")
	cmdflag := getopt.StringLong("command", 'c', "", /*TODO: Help description*/ "")
	configflag := getopt.StringLong("config", 'C', defaultconfpath, "Sets the path to Hilbish's config")
	getopt.BoolLong("login", 'l', "Makes Hilbish act like a login shell")
	getopt.BoolLong("interactive", 'i', "Force Hilbish to be an interactive shell")
	getopt.BoolLong("noexec", 'n', "Don't execute and only report Lua syntax errors")

	getopt.Parse()
	loginshflag := getopt.Lookup('l').Seen()
	interactiveflag := getopt.Lookup('i').Seen()
	noexecflag := getopt.Lookup('n').Seen()

	if *cmdflag == "" || interactiveflag {
		interactive = true
	}

	if getopt.NArgs() > 0 {
		interactive = false
	}

	if noexecflag {
		noexecute = true
	}

	// first arg, first character
	if loginshflag || os.Args[0][0] == '-' {
		login = true
	}

	if *verflag {
		fmt.Printf("Hilbish v%s\n", version)
		os.Exit(0)
	}

	// Set $SHELL if the user wants to
	if *setshflag {
		os.Setenv("SHELL", os.Args[0])
	}

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
	LuaInit()
	RunLogin()
	RunConfig(*configflag)

	readline.Completer = readline.FilenameCompleter
	readline.LoadHistory(homedir + "/.hilbish-history")

	RunInput(*cmdflag)
	if getopt.NArgs() > 0 {
		err := l.DoFile(getopt.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	for interactive {
		running = false

		input, err := readline.String(fmtPrompt())

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
		if len(input) == 0 { continue }

		if strings.HasSuffix(input, "\\") {
			for {
				input, err = ContinuePrompt(strings.TrimSuffix(input, "\\"))

				if err != nil || !strings.HasSuffix(input, "\\") {
					break
				}
			}
		}
		running = true
		RunInput(input)

		termwidth, _, err := term.GetSize(0)
		if err != nil {
			continue
		}
		fmt.Printf("\u001b[7m∆\u001b[0m" + strings.Repeat(" ", termwidth - 1) + "\r")
	}
}

func ContinuePrompt(prev string) (string, error) {
	hooks.Em.Emit("multiline", nil)
	cont, err := readline.String(multilinePrompt)
	if err != nil {
		fmt.Println("")
		return "", err
	}
	cont = strings.TrimSpace(cont)

	return prev + " " + strings.TrimSuffix(cont, "\n"), nil
}

// This semi cursed function formats our prompt (obviously)
func fmtPrompt() string {
	host, _ := os.Hostname()
	cwd, _ := os.Getwd()

	cwd = strings.Replace(cwd, curuser.HomeDir, "~", 1)

	args := []string{
		"d", cwd,
		"D", filepath.Base(cwd),
		"h", host,
		"u", curuser.Username,
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
			readline.ReplaceLine("", 0)
			readline.RefreshLine()
		}
	}
}
