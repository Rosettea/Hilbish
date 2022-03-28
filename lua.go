package main

import (
	"fmt"
	"os"

	"hilbish/util"
	"hilbish/golibs/bait"
	"hilbish/golibs/commander"
	"hilbish/golibs/fs"
/*
	"hilbish/golibs/terminal"
*/
	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib"
)

var minimalconf = `hilbish.prompt '& '`

func luaInit() {
	l = rt.New(os.Stdout)
	lib.LoadAll(l)

	lib.LoadLibs(l, hilbishLoader)
	// yes this is stupid, i know
	util.DoString(l, "hilbish = require 'hilbish'")

	// Add fs and terminal module module to Lua
	lib.LoadLibs(l, fs.Loader)
/*
	l.PreloadModule("terminal", terminal.Loader)
*/

	cmds := commander.New()
	// When a command from Lua is added, register it for use
	cmds.Events.On("commandRegister", func(cmdName string, cmd *rt.Closure) {
		commands[cmdName] = cmd
	})
	cmds.Events.On("commandDeregister", func(cmdName string) {
		delete(commands, cmdName)
	})
	lib.LoadLibs(l, cmds.Loader)

	hooks = bait.New()
	lib.LoadLibs(l, hooks.Loader)

	// Add Ctrl-C handler
	hooks.Em.On("signal.sigint", func() {
		if !interactive {
			os.Exit(0)
		}
	})

	// Add more paths that Lua can require from
	err := util.DoString(l, "package.path = package.path .. " + requirePaths)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not add preload paths! Libraries will be missing. This shouldn't happen.")
	}

	err = util.DoFile(l, "prelude/init.lua")
	if err != nil {
		err = util.DoFile(l, preloadPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Missing preload file, builtins may be missing.")
		}
	}
	fmt.Println(err)
}

func runConfig(confpath string) {
	if !interactive {
		return
	}
	err := util.DoFile(l, confpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "\nAn error has occured while loading your config! Falling back to minimal default config.")
		util.DoString(l, minimalconf)
	}
}
