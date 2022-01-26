// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"hilbish/util"

	"github.com/pborman/getopt"
	"github.com/yuin/gopher-lua"
	"mvdan.cc/sh/v3/interp"
)

var exports = map[string]lua.LGFunction {
	"alias": hlalias,
	"appendPath": hlappendPath,
	"cwd": hlcwd,
	"exec": hlexec,
	"flag": hlflag,
	"multiprompt": hlmlprompt,
	"prependPath": hlprependPath,
	"prompt": hlprompt,
	"interval": hlinterval,
	"read": hlread,
	"run": hlrun,
	"timeout": hltimeout,
}

var greeting string

func hilbishLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)

	host, _ := os.Hostname()
	username := curuser.Username

	greeting = `Welcome to {magenta}Hilbish{reset}, {cyan}` + curuser.Username + `{reset}.
The nice lil shell for {blue}Lua{reset} fanatics!
`

	if runtime.GOOS == "windows" {
		username = strings.Split(username, "\\")[1] // for some reason Username includes the hostname on windows
	}

	util.SetField(L, mod, "ver", lua.LString(version), "Hilbish version")
	util.SetField(L, mod, "user", lua.LString(username), "Username of user")
	util.SetField(L, mod, "host", lua.LString(host), "Host name of the machine")
	util.SetField(L, mod, "home", lua.LString(curuser.HomeDir), "Home directory of the user")
	util.SetField(L, mod, "dataDir", lua.LString(dataDir), "Directory for Hilbish's data files")
	util.SetField(L, mod, "interactive", lua.LBool(interactive), "If this is an interactive shell")
	util.SetField(L, mod, "login", lua.LBool(interactive), "Whether this is a login shell")
	util.SetField(L, mod, "greeting", lua.LString(greeting), "Hilbish's welcome message for interactive shells. It has Lunacolors formatting.")
	util.Document(L, mod, "Hilbish's core API, containing submodules and functions which relate to the shell itself.")

	// hilbish.userDir table
	hshuser := L.NewTable()
	userConfigDir, _ := os.UserConfigDir()
	userDataDir := ""
	// i honestly dont know what directories to use for this
	switch runtime.GOOS {
	case "linux":
		userDataDir = getenv("XDG_DATA_HOME", curuser.HomeDir + "/.local/share")
	default:
		userDataDir = filepath.Join(userConfigDir)
	}

	util.SetField(L, hshuser, "config", lua.LString(userConfigDir), "User's config directory")
	util.SetField(L, hshuser, "data", lua.LString(userDataDir), "XDG data directory")
	util.Document(L, hshuser, "User directories to store configs and/or modules.")
	L.SetField(mod, "userDir", hshuser)

	// hilbish.aliases table
	aliases = NewAliases()
	aliasesModule := aliases.Loader(L)
	util.Document(L, aliasesModule, "Alias inferface for Hilbish.")
	L.SetField(mod, "aliases", aliasesModule)

	L.Push(mod)

	return 1
}

// run(cmd)
// Runs `cmd` in Hilbish's sh interpreter.
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

// flag(f)
// Checks if the `f` flag has been passed to Hilbish.
func hlflag(L *lua.LState) int {
	flagchar := L.CheckString(1)

	flag := getopt.Lookup([]rune(flagchar)[0])
	if flag == nil {
		L.Push(lua.LNil)
		return 1
	}

	passed := flag.Seen()
	L.Push(lua.LBool(passed))

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
func hlread(L *lua.LState) int {
	luaprompt := L.CheckString(1)
	lualr := newLineReader(luaprompt)

	input, err := lualr.Read()
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(input))
	return 1
}

/* prompt(str)
Changes the shell prompt to `str`
There are a few verbs that can be used in the prompt text.
These will be formatted and replaced with the appropriate values.
`%d` - Current working directory
`%u` - Name of current user
`%h` - Hostname of device */
func hlprompt(L *lua.LState) int {
	prompt = L.CheckString(1)

	return 0
}

// multiprompt(str)
// Changes the continued line prompt to `str`
func hlmlprompt(L *lua.LState) int {
	multilinePrompt = L.CheckString(1)

	return 0
}

// alias(cmd, orig)
// Sets an alias of `orig` to `cmd`
func hlalias(L *lua.LState) int {
	alias := L.CheckString(1)
	source := L.CheckString(2)

	aliases.Add(alias, source)

	return 1
}

// appendPath(dir)
// Appends `dir` to $PATH
func hlappendPath(L *lua.LState) int {
	// check if dir is a table or a string
	arg := L.Get(1)
	fmt.Println(arg.Type())
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
func hlexec(L *lua.LState) int {
	cmd := L.CheckString(1)
	cmdArgs, _ := splitInput(cmd)
	cmdPath, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		fmt.Println(err)
		// if we get here, cmdPath will be nothing
		// therefore nothing will run
	}

	// syscall.Exec requires an absolute path to a binary
	// path, args, string slice of environments
	// TODO: alternative for windows
	syscall.Exec(cmdPath, cmdArgs, os.Environ())
	return 0 // random thought: does this ever return?
}

// goro(fn)
// Puts `fn` in a goroutine
func hlgoroutine(L *lua.LState) int {
	fn := L.CheckFunction(1)
	argnum := L.GetTop()
	args := make([]lua.LValue, argnum)
	for i := 1; i <= argnum; i++ {
		args[i - 1] = L.Get(i)
	}

	// call fn
	go func() {
		L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, args...)
	}()

	return 0
}

// timeout(cb, time)
// Runs the `cb` function after `time` in milliseconds
func hltimeout(L *lua.LState) int {
	cb := L.CheckFunction(1)
	ms := L.CheckInt(2)

	timeout := time.Duration(ms) * time.Millisecond
	time.Sleep(timeout)

	L.CallByParam(lua.P{
		Fn:      cb,
		NRet:    0,
		Protect: true,
	})
	return 0
}

// interval(cb, time)
// Runs the `cb` function every `time` milliseconds
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
// `cb` must be a function that returns a table of the entries to complete.
// Nested tables will be used as sub-completions.
func hlcomplete(L *lua.LState) int {
	scope := L.CheckString(1)
	cb := L.CheckFunction(2)

	luaCompletions[scope] = cb

	return 0
}

// prependPath(dir)
// Prepends `dir` to $PATH
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
