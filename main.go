package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"hilbish/golibs/bait"

	"github.com/pborman/getopt"
	"github.com/yuin/gopher-lua"
	"github.com/maxlandon/readline"
	"golang.org/x/term"
)

var (
	l *lua.LState
	lr *lineReader

	commands = map[string]*lua.LFunction{}
	luaCompletions = map[string]*lua.LFunction{}

	confDir string
	userDataDir string
	curuser *user.User

	hooks bait.Bait
	defaultConfPath string
	defaultHistPath string
)

func main() {
	curuser, _ = user.Current()
	homedir := curuser.HomeDir
	confDir, _ = os.UserConfigDir()
	preloadPath = strings.Replace(preloadPath, "~", homedir, 1)
	sampleConfPath = strings.Replace(sampleConfPath, "~", homedir, 1)

	// i honestly dont know what directories to use for this
	switch runtime.GOOS {
	case "linux":
		userDataDir = getenv("XDG_DATA_HOME", curuser.HomeDir + "/.local/share")
	default:
		// this is fine on windows, dont know about others
		userDataDir = confDir
	}

	if defaultConfDir == "" {
		// we'll add *our* default if its empty (wont be if its changed comptime)
		defaultConfDir = filepath.Join(confDir, "hilbish")
	} else {
		// else do ~ substitution
		defaultConfDir = expandHome(defaultHistDir)
	}
	defaultConfPath = filepath.Join(defaultConfDir, "init.lua")
	if defaultHistDir == "" {
		defaultHistDir = filepath.Join(userDataDir, "hilbish")
	} else {
		defaultHistDir = expandHome(defaultHistDir)
	}
	defaultHistPath = filepath.Join(defaultHistDir, ".hilbish-history")
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

	go handleSignals()
	luaInit()
	lr = newLineReader("", false)
	// If user's config doesn't exixt,
	if _, err := os.Stat(defaultConfPath); os.IsNotExist(err) && *configflag == defaultConfPath {
		// Read default from current directory
		// (this is assuming the current dir is Hilbish's git)
		_, err := os.ReadFile(".hilbishrc.lua")
		confpath := ".hilbishrc.lua"
		if err != nil {
			// If it wasnt found, go to the real sample conf
			_, err = os.ReadFile(sampleConfPath)
			confpath = sampleConfPath
			if err != nil {
				fmt.Println("could not find .hilbishrc.lua or", sampleConfPath)
				return
			}
		}

		runConfig(confpath)
	} else {
		runConfig(*configflag)
	}

	if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		for scanner.Scan() {
			text := scanner.Text()
			runInput(text, text)
		}
	}

	if *cmdflag != "" {
		runInput(*cmdflag, *cmdflag)
	}

	if getopt.NArgs() > 0 {
		luaArgs := l.NewTable()
		for _, arg := range getopt.Args() {
			luaArgs.Append(lua.LString(arg))
		}

		l.SetGlobal("args", luaArgs)
		err := l.DoFile(getopt.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	initialized = true
input:
	for interactive {
		running = false

		input, err := lr.Read()

		if err == io.EOF {
			// Exit if user presses ^D (ctrl + d)
			hooks.Em.Emit("hilbish.exit")
			break
		}
		if err != nil {
			if err != readline.CtrlC {
				// If we get a completely random error, print
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Println("^C")
			continue
		}
		oldInput := input

		input = strings.TrimSpace(input)
		if len(input) == 0 {
			running = true
			hooks.Em.Emit("command.exit", 0)
			continue
		}

		if strings.HasSuffix(input, "\\") {
			for {
				input, err = continuePrompt(input)
				if err != nil {
					running = true
					lr.SetPrompt(fmtPrompt(prompt))
					goto input // continue inside nested loop
				}
				if !strings.HasSuffix(input, "\\") {
					break
				}
			}
		}

		runInput(input, oldInput)

		termwidth, _, err := term.GetSize(0)
		if err != nil {
			continue
		}
		fmt.Printf("\u001b[7m∆\u001b[0m" + strings.Repeat(" ", termwidth - 1) + "\r")
	}
}

func continuePrompt(prev string) (string, error) {
	hooks.Em.Emit("multiline", nil)
	lr.SetPrompt(multilinePrompt)
	cont, err := lr.Read()
	if err != nil {
		return "", err
	}
	cont = strings.TrimSpace(cont)

	return prev + strings.TrimSuffix(cont, "\n"), nil
}

// This semi cursed function formats our prompt (obviously)
func fmtPrompt(prompt string) string {
	host, _ := os.Hostname()
	cwd, _ := os.Getwd()

	if strings.HasPrefix(cwd, curuser.HomeDir) {
		cwd = "~" + strings.TrimPrefix(cwd, curuser.HomeDir)
	}
	username := curuser.Username
	// this will be baked into binary since GOOS is a constant
	if runtime.GOOS == "windows" {
		username = strings.Split(username, "\\")[1] // for some reason Username includes the hostname on windows
	}

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

func handleHistory(cmd string) {
	lr.AddHistory(cmd)
	// TODO: load history again (history shared between sessions like this ye)
}

func expandHome(path string) string {
	homedir := curuser.HomeDir

	return strings.Replace(defaultHistDir, "~", homedir, 1)
}

func removeDupes(slice []string) []string {
	all := make(map[string]bool)
	newSlice := []string{}
	for _, item := range slice {
		if _, val := all[item]; !val {
			all[item] = true
			newSlice = append(newSlice, item)
		}
	}

	return newSlice
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
