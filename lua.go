package main

import (
	"fmt"
	"os"

	"hilbish/golibs/bait"
	"hilbish/golibs/commander"
	"hilbish/golibs/fs"
	"hilbish/golibs/terminal"

	"github.com/yuin/gopher-lua"
)

var minimalconf = `hilbish.prompt '& '`

func luaInit() {
	l = lua.NewState()
	l.OpenLibs()

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

	// Add more paths that Lua can require from
	l.DoString("package.path = package.path .. " + requirePaths)

	err := l.DoFile("prelude/init.lua")
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
