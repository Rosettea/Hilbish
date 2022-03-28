package main

import (
	"fmt"
	"os"

	"hilbish/util"
	"hilbish/golibs/bait"
/*
	"hilbish/golibs/commander"
	"hilbish/golibs/fs"
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
	chunk, _ := l.CompileAndLoadLuaChunk("", []byte("hilbish = require 'hilbish'"), rt.TableValue(l.GlobalEnv()))
	_, err := rt.Call1(l.MainThread(), rt.FunctionValue(chunk))
	fmt.Println("hsh load", err)

	// Add fs and terminal module module to Lua
/*	l.PreloadModule("fs", fs.Loader)
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
*/

	hooks = bait.New()
//	l.PreloadModule("bait", hooks.Loader)
	// Add Ctrl-C handler
	hooks.Em.On("signal.sigint", func() {
		if !interactive {
			os.Exit(0)
		}
	})

	// Add more paths that Lua can require from
	chunk, _ = l.CompileAndLoadLuaChunk("", []byte("package.path = package.path .. " + requirePaths), rt.TableValue(l.GlobalEnv()))
	_, err = rt.Call1(l.MainThread(), rt.FunctionValue(chunk))
	fmt.Println("package path", err)

	data, err := os.ReadFile("prelude/init.lua")
	if err != nil {
		data, err = os.ReadFile(preloadPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Missing preload file, builtins may be missing.")
		}
	}
	chunk, _ = l.CompileAndLoadLuaChunk("", data, rt.TableValue(l.GlobalEnv()))
	_, err = rt.Call1(l.MainThread(), rt.FunctionValue(chunk))
	fmt.Println("prelude", err)
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
