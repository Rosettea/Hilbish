package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"os/signal"
	"os/user"
	"path/filepath"
	"time"

	"hilbish/golibs/bait"

	"github.com/pborman/getopt"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"mvdan.cc/sh/v3/interp"
	"golang.org/x/term"
)

var (
	l *lua.LState
	lr *LineReader

	commands = map[string]*lua.LFunction{}
	aliases = map[string]string{}

	homedir string
	curuser *user.User

	hooks bait.Bait
	defaultConfPath string

	runner *interp.Runner	
)

func main() {
	homedir, _ = os.UserHomeDir()
	curuser, _ = user.Current()
	preloadPath = strings.Replace(preloadPath, "~", homedir, 1)
	sampleConfPath = strings.Replace(sampleConfPath, "~", homedir, 1)

	if defaultConfDir == "" {
		// we'll add *our* default if its empty (wont be if its changed comptime)
		defaultConfPath = filepath.Join(homedir, "/.hilbishrc.lua")
	} else {
		// else do ~ substitution
		defaultConfPath = filepath.Join(strings.Replace(defaultConfDir, "~", homedir, 1), ".hilbishrc.lua")
	}

	helpflag := getopt.BoolLong("help", 'h', "Prints Hilbish flags")
	verflag := getopt.BoolLong("version", 'v', "Prints Hilbish version")
	setshflag := getopt.BoolLong("setshellenv", 'S', "Sets $SHELL to Hilbish's executed path")
	cmdflag := getopt.StringLong("command", 'c', "", "Executes a command on startup")
	configflag := getopt.StringLong("config", 'C', defaultConfPath, "Sets the path to Hilbish's config")
	getopt.BoolLong("login", 'l', "Force Hilbish to be a login shell")
	getopt.BoolLong("interactive", 'i', "Force Hilbish to be an interactive shell")
	getopt.BoolLong("noexec", 'n', "Don't execute and only report Lua syntax errors")

	getopt.Parse()
	loginshflag := getopt.Lookup('l').Seen()
	interactiveflag := getopt.Lookup('i').Seen()
	noexecflag := getopt.Lookup('n').Seen()

	if *helpflag {
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	}

	if *cmdflag == "" || interactiveflag {
		interactive = true
	}

	if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		interactive = false
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
		fmt.Printf("Hilbish %s\n", version)
		os.Exit(0)
	}

	// Set $SHELL if the user wants to
	if *setshflag {
		os.Setenv("SHELL", os.Args[0])
	}

	// If user's config doesn't exixt,
	if _, err := os.Stat(defaultConfPath); os.IsNotExist(err) {
		// Read default from current directory
		// (this is assuming the current dir is Hilbish's git)
		input, err := os.ReadFile(".hilbishrc.lua")
		if err != nil {
			// If it wasnt found, go to the real sample conf
			input, err = os.ReadFile(sampleConfPath)
			if err != nil {
				fmt.Println("could not find .hilbishrc.lua or", sampleConfPath)
				return
			}
		}

		// Create it using either default config we found
		err = os.WriteFile(defaultConfPath, input, 0644)
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

	lr = NewLineReader("")

	if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		for scanner.Scan() {
			RunInput(scanner.Text())
		}
	}

	if *cmdflag != "" {
		RunInput(*cmdflag)
	}

	if getopt.NArgs() > 0 {
		l.SetGlobal("args", luar.New(l, getopt.Args()))
		err := l.DoFile(getopt.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Execute handler for sh runner
	exechandle := func(ctx context.Context, args []string) error {
		hc := interp.HandlerCtx(ctx)
		_, argstring := splitInput(strings.Join(args, " "))

		// If alias was found, use command alias
		if aliases[args[0]] != "" {
			alias := aliases[args[0]]
			argstring = alias + strings.TrimPrefix(argstring, args[0])
			cmdArgs, _ := splitInput(argstring)
			args = cmdArgs
		}

		// If command is defined in Lua then run it
		if commands[args[0]] != nil {
			err := l.CallByParam(lua.P{
				Fn: commands[args[0]],
				NRet:    1,
				Protect: true,
			}, luar.New(l, args[1:]))
			luaexitcode := l.Get(-1)
			var exitcode uint8 = 0

			l.Pop(1)

			if code, ok := luaexitcode.(lua.LNumber); luaexitcode != lua.LNil && ok {
				exitcode = uint8(code)
			}

			if err != nil {
				fmt.Fprintln(os.Stderr,
					"Error in command:\n\n" + err.Error())
			}
			hooks.Em.Emit("command.exit", exitcode)
			return interp.NewExitStatus(exitcode)
		}

		if _, err := interp.LookPathDir(hc.Dir, hc.Env, args[0]); err != nil {
			hooks.Em.Emit("command.not-found", args[0])
			return interp.NewExitStatus(127)
		}

		return interp.DefaultExecHandler(2 * time.Second)(ctx, args)
	}
	// Setup sh runner outside of input label
	runner, _ = interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.ExecHandler(exechandle),
	)

input:
	for interactive {
		running = false

		lr.SetPrompt(fmtPrompt())
		input, err := lr.Read()

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
				if err != nil {
					goto input // continue inside nested loop
				}
				if !strings.HasSuffix(input, "\\") {
					break
				}
			}
		}
		running = true
		HandleHistory(input)
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
	lr.SetPrompt(multilinePrompt)
	cont, err := lr.Read()
	if err != nil {
		fmt.Println("")
		return "", err
	}
	cont = strings.TrimSpace(cont)

	return prev + strings.TrimSuffix(cont, "\n"), nil
}

// This semi cursed function formats our prompt (obviously)
func fmtPrompt() string {
	host, _ := os.Hostname()
	cwd, _ := os.Getwd()

	cwd = strings.Replace(cwd, curuser.HomeDir, "~", 1)
	username := strings.Split(curuser.Username, "\\")[1] // for some reason Username includes the hostname on windows

	args := []string{
		"d", cwd,
		"D", filepath.Base(cwd),
		"h", host,
		"u", username,
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
			lr.ClearInput()
		}
	}
}

func HandleHistory(cmd string) {
	lr.AddHistory(cmd)
	// TODO: load history again (history shared between sessions like this ye)
}

