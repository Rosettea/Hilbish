// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"hilbish/util"

	"github.com/yuin/gopher-lua"
	"github.com/maxlandon/readline"
	"github.com/blackfireio/osinfo"
	"mvdan.cc/sh/v3/interp"
)

var exports = map[string]lua.LGFunction {
	"alias": hlalias,
	"appendPath": hlappendPath,
	"complete": hlcomplete,
	"cwd": hlcwd,
	"exec": hlexec,
	"runnerMode": hlrunnerMode,
	"goro": hlgoro,
	"multiprompt": hlmlprompt,
	"prependPath": hlprependPath,
	"prompt": hlprompt,
	"inputMode": hlinputMode,
	"interval": hlinterval,
	"read": hlread,
	"run": hlrun,
	"timeout": hltimeout,
	"which": hlwhich,
}

var greeting string
var hshMod *lua.LTable

func hilbishLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	hshMod = mod

	host, _ := os.Hostname()
	username := curuser.Username

	if runtime.GOOS == "windows" {
		username = strings.Split(username, "\\")[1] // for some reason Username includes the hostname on windows
	}

	greeting = `Welcome to {magenta}Hilbish{reset}, {cyan}` + username + `{reset}.
The nice lil shell for {blue}Lua{reset} fanatics!
Check out the {blue}{bold}guide{reset} command to get started.
`

	util.SetField(L, mod, "ver", lua.LString(version), "Hilbish version")
	util.SetField(L, mod, "user", lua.LString(username), "Username of user")
	util.SetField(L, mod, "host", lua.LString(host), "Host name of the machine")
	util.SetField(L, mod, "home", lua.LString(curuser.HomeDir), "Home directory of the user")
	util.SetField(L, mod, "dataDir", lua.LString(dataDir), "Directory for Hilbish's data files")
	util.SetField(L, mod, "interactive", lua.LBool(interactive), "If this is an interactive shell")
	util.SetField(L, mod, "login", lua.LBool(interactive), "Whether this is a login shell")
	util.SetField(L, mod, "greeting", lua.LString(greeting), "Hilbish's welcome message for interactive shells. It has Lunacolors formatting.")
	util.SetField(l, mod, "vimMode", lua.LNil, "Current Vim mode of Hilbish (nil if not in Vim mode)")
	util.SetField(l, hshMod, "exitCode", lua.LNumber(0), "Exit code of last exected command")
	util.Document(L, mod, "Hilbish's core API, containing submodules and functions which relate to the shell itself.")

	// hilbish.userDir table
	hshuser := L.NewTable()

	util.SetField(L, hshuser, "config", lua.LString(confDir), "User's config directory")
	util.SetField(L, hshuser, "data", lua.LString(userDataDir), "XDG data directory")
	util.Document(L, hshuser, "User directories to store configs and/or modules.")
	L.SetField(mod, "userDir", hshuser)

	// hilbish.os table
	hshos := L.NewTable()
	info, _ := osinfo.GetOSInfo()

	util.SetField(L, hshos, "family", lua.LString(info.Family), "Family name of the current OS")
	util.SetField(L, hshos, "name", lua.LString(info.Name), "Pretty name of the current OS")
	util.SetField(L, hshos, "version", lua.LString(info.Version), "Version of the current OS")
	util.Document(L, hshos, "OS info interface")
	L.SetField(mod, "os", hshos)

	// hilbish.aliases table
	aliases = newAliases()
	aliasesModule := aliases.Loader(L)
	util.Document(L, aliasesModule, "Alias inferface for Hilbish.")
	L.SetField(mod, "aliases", aliasesModule)

	// hilbish.history table
	historyModule := lr.Loader(L)
	util.Document(L, historyModule, "History interface for Hilbish.")
	L.SetField(mod, "history", historyModule)

	// hilbish.completions table
	hshcomp := L.NewTable()

	util.SetField(L, hshcomp, "files", L.NewFunction(luaFileComplete), "Completer for files")
	util.SetField(L, hshcomp, "bins", L.NewFunction(luaBinaryComplete), "Completer for executables/binaries")
	util.Document(L, hshcomp, "Completions interface for Hilbish.")
	L.SetField(mod, "completion", hshcomp)

	// hilbish.runner table
	runnerModule := runnerModeLoader(L)
	util.Document(L, runnerModule, "Runner/exec interface for Hilbish.")
	L.SetField(mod, "runner", runnerModule)

	// hilbish.jobs table
	jobs = newJobHandler()
	jobModule := jobs.loader(L)
	util.Document(L, jobModule, "(Background) job interface.")
	L.SetField(mod, "jobs", jobModule)

	L.Push(mod)

	return 1
}

func luaFileComplete(L *lua.LState) int {
	query := L.CheckString(1)
	ctx := L.CheckString(2)
	fields := L.CheckTable(3)

	var fds []string
	fields.ForEach(func(k lua.LValue, v lua.LValue) {
		fds = append(fds, v.String())
	})

	completions := fileComplete(query, ctx, fds)
	luaComps := L.NewTable()

	for _, comp := range completions {
		luaComps.Append(lua.LString(comp))
	}

	L.Push(luaComps)

	return 1
}

func luaBinaryComplete(L *lua.LState) int {
	query := L.CheckString(1)
	ctx := L.CheckString(2)
	fields := L.CheckTable(3)

	var fds []string
	fields.ForEach(func(k lua.LValue, v lua.LValue) {
		fds = append(fds, v.String())
	})

	completions, _ := binaryComplete(query, ctx, fds)
	luaComps := L.NewTable()

	for _, comp := range completions {
		luaComps.Append(lua.LString(comp))
	}

	L.Push(luaComps)

	return 1
}

func setVimMode(mode string) {
	util.SetField(l, hshMod, "vimMode", lua.LString(mode), "Current Vim mode of Hilbish (nil if not in Vim mode)")
	hooks.Em.Emit("hilbish.vimMode", mode)
}

func unsetVimMode() {
	util.SetField(l, hshMod, "vimMode", lua.LNil, "Current Vim mode of Hilbish (nil if not in Vim mode)")
}

// run(cmd)
// Runs `cmd` in Hilbish's sh interpreter.
// --- @param cmd string
func hlrun(L *lua.LState) int {
	var exitcode uint8
	cmd := L.CheckString(1)
	err := execCommand(cmd)

	if code, ok := interp.IsExitStatus(err); ok {
		exitcode = code
	} else if err != nil {
		exitcode = 1
	}

	L.Push(lua.LNumber(exitcode))
	return 1
}

// cwd()
// Returns the current directory of the shell
func hlcwd(L *lua.LState) int {
	cwd, _ := os.Getwd()

	L.Push(lua.LString(cwd))

	return 1
}

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

// read(prompt) -> input?
// Read input from the user, using Hilbish's line editor/input reader.
// This is a separate instance from the one Hilbish actually uses.
// Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)
// --- @param prompt string
func hlread(L *lua.LState) int {
	luaprompt := L.CheckString(1)
	lualr := newLineReader("", true)
	lualr.SetPrompt(luaprompt)

	input, err := lualr.Read()
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(input))
	return 1
}

/*
prompt(str)
Changes the shell prompt to `str`
There are a few verbs that can be used in the prompt text.
These will be formatted and replaced with the appropriate values.
`%d` - Current working directory
`%u` - Name of current user
`%h` - Hostname of device
--- @param str string
*/
func hlprompt(L *lua.LState) int {
	prompt = L.CheckString(1)
	lr.SetPrompt(fmtPrompt(prompt))

	return 0
}

// multiprompt(str)
// Changes the continued line prompt to `str`
// --- @param str string
func hlmlprompt(L *lua.LState) int {
	multilinePrompt = L.CheckString(1)

	return 0
}

// alias(cmd, orig)
// Sets an alias of `cmd` to `orig`
// --- @param cmd string
// --- @param orig string
func hlalias(L *lua.LState) int {
	alias := L.CheckString(1)
	source := L.CheckString(2)

	aliases.Add(alias, source)

	return 1
}

// appendPath(dir)
// Appends `dir` to $PATH
// --- @param dir string|table
func hlappendPath(L *lua.LState) int {
	// check if dir is a table or a string
	arg := L.Get(1)
	if arg.Type() == lua.LTTable {
		arg.(*lua.LTable).ForEach(func(k lua.LValue, v lua.LValue) {
			appendPath(v.String())
		})
	} else if arg.Type() == lua.LTString {
		appendPath(arg.String())
	} else {
		L.RaiseError("bad argument to appendPath (expected string or table, got %v)", L.Get(1).Type().String())
	}

	return 0
}

func appendPath(dir string) {
	dir = strings.Replace(dir, "~", curuser.HomeDir, 1)
	pathenv := os.Getenv("PATH")

	// if dir isnt already in $PATH, add it
	if !strings.Contains(pathenv, dir) {
		os.Setenv("PATH", pathenv + string(os.PathListSeparator) + dir)
	}
}

// exec(cmd)
// Replaces running hilbish with `cmd`
// --- @param cmd string
func hlexec(L *lua.LState) int {
	cmd := L.CheckString(1)
	cmdArgs, _ := splitInput(cmd)
	if runtime.GOOS != "windows" {
		cmdPath, err := exec.LookPath(cmdArgs[0])
		if err != nil {
			fmt.Println(err)
			// if we get here, cmdPath will be nothing
			// therefore nothing will run
		}

		// syscall.Exec requires an absolute path to a binary
		// path, args, string slice of environments
		syscall.Exec(cmdPath, cmdArgs, os.Environ())
	} else {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
		os.Exit(0)
	}

	return 0
}

// goro(fn)
// Puts `fn` in a goroutine
// --- @param fn function
func hlgoro(L *lua.LState) int {
	fn := L.CheckFunction(1)
	argnum := L.GetTop()
	args := make([]lua.LValue, argnum)
	for i := 1; i <= argnum; i++ {
		args[i - 1] = L.Get(i)
	}

	// call fn
	go func() {
		if err := L.CallByParam(lua.P{
			Fn: fn,
			NRet: 0,
			Protect: true,
		}, args...); err != nil {
			fmt.Fprintln(os.Stderr, "Error in goro function:\n\n", err)
		}
	}()

	return 0
}

// timeout(cb, time)
// Runs the `cb` function after `time` in milliseconds
// --- @param cb function
// --- @param time number
func hltimeout(L *lua.LState) int {
	cb := L.CheckFunction(1)
	ms := L.CheckInt(2)

	timeout := time.Duration(ms) * time.Millisecond
	time.Sleep(timeout)

	if err := L.CallByParam(lua.P{
		Fn: cb,
		NRet: 0,
		Protect: true,
	}); err != nil {
		fmt.Fprintln(os.Stderr, "Error in goro function:\n\n", err)
	}
	return 0
}

// interval(cb, time)
// Runs the `cb` function every `time` milliseconds
// --- @param cb function
// --- @param time number
func hlinterval(L *lua.LState) int {
	intervalfunc := L.CheckFunction(1)
	ms := L.CheckInt(2)
	interval := time.Duration(ms) * time.Millisecond

	ticker := time.NewTicker(interval)
	stop := make(chan lua.LValue)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := L.CallByParam(lua.P{
					Fn: intervalfunc,
					NRet: 0,
					Protect: true,
				}); err != nil {
					fmt.Fprintln(os.Stderr, "Error in interval function:\n\n", err)
					stop <- lua.LTrue // stop the interval
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	L.Push(lua.LChannel(stop))
	return 1
}

// complete(scope, cb)
// Registers a completion handler for `scope`.
// A `scope` is currently only expected to be `command.<cmd>`,
// replacing <cmd> with the name of the command (for example `command.git`).
// `cb` must be a function that returns a table of "completion groups."
// A completion group is a table with the keys `items` and `type`.
// `items` being a table of items and `type` being the display type of
// `grid` (the normal file completion display) or `list` (with a description)
// --- @param scope string
// --- @param cb function
func hlcomplete(L *lua.LState) int {
	scope := L.CheckString(1)
	cb := L.CheckFunction(2)

	luaCompletions[scope] = cb

	return 0
}

// prependPath(dir)
// Prepends `dir` to $PATH
// --- @param dir string
func hlprependPath(L *lua.LState) int {
	dir := L.CheckString(1)
	dir = strings.Replace(dir, "~", curuser.HomeDir, 1)
	pathenv := os.Getenv("PATH")

	// if dir isnt already in $PATH, add in
	if !strings.Contains(pathenv, dir) {
		os.Setenv("PATH", dir + string(os.PathListSeparator) + pathenv)
	}

	return 0
}

// which(binName)
// Searches for an executable called `binName` in the directories of $PATH
// --- @param binName string
func hlwhich(L *lua.LState) int {
	binName := L.CheckString(1)
	path, err := exec.LookPath(binName)
	if err != nil {
		l.Push(lua.LNil)
		return 1
	}

	l.Push(lua.LString(path))
	return 1
}

// inputMode(mode)
// Sets the input mode for Hilbish's line reader. Accepts either emacs for vim
// --- @param mode string
func hlinputMode(L *lua.LState) int {
	mode := L.CheckString(1)
	switch mode {
		case "emacs":
			unsetVimMode()
			lr.rl.InputMode = readline.Emacs
		case "vim":
			setVimMode("insert")
			lr.rl.InputMode = readline.Vim
		default: L.RaiseError("inputMode: expected vim or emacs, received " + mode)
	}
	return 0
}

// runnerMode(mode)
// Sets the execution/runner mode for interactive Hilbish. This determines whether
// Hilbish wll try to run input as Lua and/or sh or only do one of either.
// Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
// sh, and lua. It also accepts a function, to which if it is passed one
// will call it to execute user input instead.
// --- @param mode string|function
func hlrunnerMode(L *lua.LState) int {
	mode := L.CheckAny(1)
	switch mode.Type() {
		case lua.LTString:
			switch mode.String() {
				// no fallthrough doesnt work so eh
				case "hybrid": fallthrough
				case "hybridRev": fallthrough
				case "lua": fallthrough
				case "sh":
					runnerMode = mode
				default: L.RaiseError("execMode: expected either a function or hybrid, hybridRev, lua, sh. Received %v", mode)
			}
		case lua.LTFunction: runnerMode = mode
		default: L.RaiseError("execMode: expected either a function or hybrid, hybridRev, lua, sh. Received %v", mode)
	}

	return 0
}
