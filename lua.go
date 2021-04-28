package main

import (
	"fmt"
	hooks "hilbish/golibs/bait"
	cmds "hilbish/golibs/commander"
	lfs "hilbish/golibs/fs"
	"os"

	lua "github.com/yuin/gopher-lua"
)

var minimalconf = `
ansikit = require 'ansikit'
prompt(ansikit.format(
		'{blue}%u {cyan}%d {green}∆{reset} '
))
`

func LuaInit(confpath string) {
	l = lua.NewState()

	l.OpenLibs()

	l.SetGlobal("_ver", lua.LString(version))

	l.SetGlobal("prompt", l.NewFunction(hshprompt))
	l.SetGlobal("multiprompt", l.NewFunction(hshmlprompt))
	l.SetGlobal("alias", l.NewFunction(hshalias))

	// Add fs module to Lua
	l.PreloadModule("fs", lfs.Loader)

	commander := cmds.New()
	// When a command from Lua is added, register it for use
	commander.Events.On("commandRegister",
		func(cmdName string, cmd *lua.LFunction) {
			commands[cmdName] = true
			l.SetField(
				l.GetTable(l.GetGlobal("commanding"),
					lua.LString("__commands")),
				cmdName,
				cmd)
		})

	l.PreloadModule("commander", commander.Loader)

	bait = hooks.New()
	l.PreloadModule("bait", bait.Loader)

	// Add more paths that Lua can require from
	l.DoString(`package.path = package.path
	.. ';./libs/?/init.lua;./?/init.lua;./?/?.lua'
	.. ';/usr/share/hilbish/libs/?/init.lua;'
	.. os.getenv 'HOME' .. '/.local/share/hilbish/libs/?/init.lua;'
	.. os.getenv 'HOME'	.. '/.local/share/hilbish/libs/?/?.lua;'
	.. os.getenv 'HOME'	.. '/.local/share/hilbish/libs/?.lua'
	`)

	err := l.DoFile("/usr/share/hilbish/preload.lua")
	if err != nil {
		err = l.DoFile("preload.lua")
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Missing preload file, builtins may be missing.")
		}
	}

	// Run config
	if !interactive {
		return
	}
	err = l.DoFile(confpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err,
			"\nAn error has occured while loading your config! Falling back to minimal default config.\n")

		l.DoString(minimalconf)
	}
}

func hshprompt(L *lua.LState) int {
	prompt = L.ToString(1)

	return 0
}

func hshmlprompt(L *lua.LState) int {
	multilinePrompt = L.ToString(1)

	return 0
}

func hshalias(L *lua.LState) int {
	alias := L.ToString(1)
	source := L.ToString(2)

	aliases[alias] = source

	return 1
}
