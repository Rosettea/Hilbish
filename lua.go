package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"hilbish/golibs/bait"
	"hilbish/golibs/commander"
	"hilbish/golibs/fs"
	"hilbish/golibs/terminal"

	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

var minimalconf = `prompt '& '`

func luaInit() {
	l = lua.NewState()
	l.OpenLibs()

	l.SetGlobal("prompt", l.NewFunction(hshprompt))
	l.SetGlobal("multiprompt", l.NewFunction(hshmlprompt))
	l.SetGlobal("alias", l.NewFunction(hshalias))
	l.SetGlobal("appendPath", l.NewFunction(hshappendPath))
	l.SetGlobal("prependPath", l.NewFunction(hshprependPath))
	l.SetGlobal("exec", l.NewFunction(hshexec))
	l.SetGlobal("goro", luar.New(l, hshgoroutine))
	l.SetGlobal("timeout", luar.New(l, hshtimeout))
	l.SetGlobal("interval", l.NewFunction(hshinterval))

	// yes this is stupid, i know
	l.PreloadModule("hilbish", hilbishLoader)
	l.DoString("hilbish = require 'hilbish'")

	// Add fs and terminal module module to Lua
	l.PreloadModule("fs", fs.Loader)
	l.PreloadModule("terminal", terminal.Loader)

	cmds := commander.New()
	// When a command from Lua is added, register it for use
	cmds.Events.On("commandRegister", func(cmdName string, cmd *lua.LFunction) {
		commands[cmdName] = cmd
	})
	cmds.Events.On("commandDeregister", func(cmdName string) {
		delete(commands, cmdName)
	})
	l.PreloadModule("commander", cmds.Loader)

	hooks = bait.New()
	l.PreloadModule("bait", hooks.Loader)

	// Add Ctrl-C handler
	hooks.Em.On("signal.sigint", func() {
		if !interactive {
			os.Exit(0)
		}
	})

	l.SetGlobal("complete", l.NewFunction(hshcomplete))

	// Add more paths that Lua can require from
	l.DoString("package.path = package.path .. " + requirePaths)

	err := l.DoFile("preload.lua")
	if err != nil {
		err = l.DoFile(preloadPath)
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Missing preload file, builtins may be missing.")
		}
	}
}
func runConfig(confpath string) {
	if !interactive {
		return
	}
	err := l.DoFile(confpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err,
			"\nAn error has occured while loading your config! Falling back to minimal default config.")

		l.DoString(minimalconf)
	}
}

func runLogin() {
	if _, err := os.Stat(curuser.HomeDir + "/.hprofile.lua"); os.IsNotExist(err) {
		return
	}
	if !login {
		return
	}
	err := l.DoFile(curuser.HomeDir + "/.hprofile.lua")
	if err != nil {
		fmt.Fprintln(os.Stderr, err,
			"\nAn error has occured while loading your login config!n")
	}
}

/* prompt(str)
Changes the shell prompt to `str`
There are a few verbs that can be used in the prompt text.
These will be formatted and replaced with the appropriate values.
`%d` - Current working directory
`%u` - Name of current user
`%h` - Hostname of device */
func hshprompt(L *lua.LState) int {
	prompt = L.CheckString(1)

	return 0
}

// multiprompt(str)
// Changes the continued line prompt to `str`
func hshmlprompt(L *lua.LState) int {
	multilinePrompt = L.CheckString(1)

	return 0
}

// alias(cmd, orig)
// Sets an alias of `orig` to `cmd`
func hshalias(L *lua.LState) int {
	alias := L.CheckString(1)
	source := L.CheckString(2)

	aliases.Add(alias, source)

	return 1
}

// appendPath(dir)
// Appends `dir` to $PATH
func hshappendPath(L *lua.LState) int {
	dir := L.CheckString(1)
	dir = strings.Replace(dir, "~", curuser.HomeDir, 1)
	pathenv := os.Getenv("PATH")

	// if dir isnt already in $PATH, add it
	if !strings.Contains(pathenv, dir) {
		os.Setenv("PATH", pathenv + string(os.PathListSeparator) + dir)
	}

	return 0
}

// exec(cmd)
// Replaces running hilbish with `cmd`
func hshexec(L *lua.LState) int {
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
func hshgoroutine(gofunc func()) {
	go gofunc()
}

// timeout(cb, time)
// Runs the `cb` function after `time` in milliseconds
func hshtimeout(timeoutfunc func(), ms int) {
	timeout := time.Duration(ms) * time.Millisecond
	time.Sleep(timeout)
	timeoutfunc()
}

// interval(cb, time)
// Runs the `cb` function every `time` milliseconds
func hshinterval(L *lua.LState) int {
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
func hshcomplete(L *lua.LState) int {
	scope := L.CheckString(1)
	cb := L.CheckFunction(2)

	luaCompletions[scope] = cb

	return 0
}

// prependPath(dir)
// Prepends `dir` to $PATH
func hshprependPath(L *lua.LState) int {
	dir := L.CheckString(1)
	dir = strings.Replace(dir, "~", curuser.HomeDir, 1)
	pathenv := os.Getenv("PATH")

	// if dir isnt already in $PATH, add in
	if !strings.Contains(pathenv, dir) {
		os.Setenv("PATH", dir + string(os.PathListSeparator) + pathenv)
	}

	return 0
}
