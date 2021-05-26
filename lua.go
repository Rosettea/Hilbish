package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"hilbish/golibs/bait"
	"hilbish/golibs/commander"
	"hilbish/golibs/fs"

	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

var minimalconf = `
ansikit = require 'ansikit'
prompt(ansikit.format(
		'{blue}%u {cyan}%d {green}∆{reset} '
))
`

func LuaInit() {
	l = lua.NewState()
	l.OpenLibs()

	l.SetGlobal("prompt", l.NewFunction(hshprompt))
	l.SetGlobal("multiprompt", l.NewFunction(hshmlprompt))
	l.SetGlobal("alias", l.NewFunction(hshalias))
	l.SetGlobal("appendPath", l.NewFunction(hshappendPath))
	l.SetGlobal("exec", l.NewFunction(hshexec))
	l.SetGlobal("goro", luar.New(l, hshgoroutine))

	// yes this is stupid, i know
	l.PreloadModule("hilbish", HilbishLoader)
	l.DoString("hilbish = require 'hilbish'")

	// Add fs module to Lua
	l.PreloadModule("fs", fs.Loader)

	cmds := commander.New()
	// When a command from Lua is added, register it for use
	cmds.Events.On("commandRegister",
		func(cmdName string, cmd *lua.LFunction) {
			commands[cmdName] = true
			l.SetField(
				l.GetTable(l.GetGlobal("commanding"),
					lua.LString("__commands")),
				cmdName,
				cmd)
		})

	l.PreloadModule("commander", cmds.Loader)

	hooks = bait.New()
	l.PreloadModule("bait", hooks.Loader)

	// Add more paths that Lua can require from
	l.DoString(`package.path = package.path
	.. ';./libs/?/init.lua;./?/init.lua;./?/?.lua'
	.. ';/usr/share/hilbish/libs/?/init.lua;'
	.. ';/usr/share/hilbish/libs/?/?.lua;'
	.. os.getenv 'HOME' .. '/.local/share/hilbish/libs/?/init.lua;'
	.. os.getenv 'HOME'	.. '/.local/share/hilbish/libs/?/?.lua;'
	.. os.getenv 'HOME'	.. '/.local/share/hilbish/libs/?.lua'
	`)

	err := l.DoFile("preload.lua")
	if err != nil {
		err = l.DoFile("/usr/share/hilbish/preload.lua")
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Missing preload file, builtins may be missing.")
		}
	}
}
func RunConfig(confpath string) {
	if !interactive {
		return
	}
	err := l.DoFile(confpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err,
			"\nAn error has occured while loading your config! Falling back to minimal default config.\n")

		l.DoString(minimalconf)
	}
}

func RunLogin() {
	if _, err := os.Stat(homedir + "/.hprofile.lua"); os.IsNotExist(err) {
		return
	}
	if !login {
		return
	}
	err := l.DoFile(homedir + "/.hprofile.lua")
	if err != nil {
		fmt.Fprintln(os.Stderr, err,
			"\nAn error has occured while loading your login config!n")
	}
}

func hshprompt(L *lua.LState) int {
	prompt = L.CheckString(1)

	return 0
}

func hshmlprompt(L *lua.LState) int {
	multilinePrompt = L.CheckString(1)

	return 0
}

func hshalias(L *lua.LState) int {
	alias := L.CheckString(1)
	source := L.CheckString(2)

	aliases[alias] = source

	return 1
}

func hshappendPath(L *lua.LState) int {
	dir := L.CheckString(1)
	dir = strings.Replace(dir, "~", curuser.HomeDir, 1)
	pathenv := os.Getenv("PATH")

	// if dir isnt already in $PATH, add it
	if !strings.Contains(pathenv, dir) {
		os.Setenv("PATH", pathenv + ":" + dir)
	}

	return 0
}

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

func hshgoroutine(gofunc func()) {
	go gofunc()
}
