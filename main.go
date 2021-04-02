package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"syscall"
	"os/signal"
	"strings"
	"io"
	"context"
	hooks "hilbish/golibs/bait"

	"github.com/akamensky/argparse"
	"github.com/bobappleyard/readline"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

)

const version = "0.2.0"
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
		// First try to run user input in Lua
		err = l.DoString(cmdString)

		if err == nil {
			// If it succeeds, add to history and prompt again
			readline.AddHistory(cmdString)
			readline.SaveHistory(homedir + "/.hilbish-history")
			bait.Em.Emit("command.success", nil)
			continue
		}

		// Split up the input
		cmdArgs, cmdString := splitInput(cmdString)
		// If there's actually no input, prompt again
		if len(cmdArgs) == 0 { continue }

		// If alias was found, use command alias
		if aliases[cmdArgs[0]] != "" {
			cmdString = aliases[cmdArgs[0]] + strings.Trim(cmdString, cmdArgs[0])
			execCommand(cmdString)
			continue
		}

		// If command is defined in Lua then run it
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
			readline.SaveHistory(homedir + "/.hilbish-history")
			continue
		}

		// Last option: use sh interpreter
		switch cmdArgs[0] {
		case "exit":
			os.Exit(0)
		default:
			err := execCommand(cmdString)
			if err != nil {
				// If input is incomplete, start multiline prompting
				if syntax.IsIncomplete(err) {
					sb := &strings.Builder{}
					for {
						done := StartMultiline(cmdString, sb)
						if done {
							break
						}
					}
				} else {
					if code, ok := interp.IsExitStatus(err); ok {
						if code > 0 {
							bait.Em.Emit("command.fail", code)
						}
					}
					fmt.Fprintln(os.Stderr, err)
				}
			} else {
				bait.Em.Emit("command.success", nil)
			}
		}
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

// do i even have to say
func HandleSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()
}

