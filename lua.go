package main

import (
	"fmt"
	"os"
	"path/filepath"

	//"hilbish/util"
	"hilbish/golibs/bait"
	"hilbish/golibs/commander"
	"hilbish/golibs/fs"
	"hilbish/golibs/snail"
	"hilbish/golibs/terminal"

	"hilbish/moonlight"
)

var minimalconf = `hilbish.prompt '& '`

func luaInit() {
	l = moonlight.NewRuntime()

	l.LoadLibrary(hilbishLoader, "hilbish")
	// yes this is stupid, i know
	l.DoString("hilbish = require 'hilbish'")

	hooks = bait.New(l)
	hooks.SetRecoverer(func(event string, handler *bait.Listener, err interface{}) {
		fmt.Println("Error in `error` hook handler:", err)
		hooks.Off(event, handler)
	})
	l.LoadLibrary(hooks.Loader, "bait")

	// Add Ctrl-C handler
	/*
		hooks.On("signal.sigint", func(...interface{}) rt.Value {
			if !interactive {
				os.Exit(0)
			}
			return rt.NilValue
		})

		lr.rl.RawInputCallback = func(r []rune) {
			hooks.Emit("hilbish.rawInput", string(r))
		}
	*/

	l.LoadLibrary(fs.Loader, "fs")
	l.LoadLibrary(terminal.Loader, "terminal")
	l.LoadLibrary(snail.Loader, "snail")

	cmds = commander.New(l)
	l.LoadLibrary(cmds.Loader, "commander")

	// Add more paths that Lua can require from
	_, err := l.DoString("package.path = package.path .. " + requirePaths)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Could not add Hilbish require paths! Libraries will be missing. This shouldn't happen.")
	}

	err1 := l.DoFile("nature/init.lua")
	if err1 != nil {
		err2 := l.DoFile(filepath.Join(dataDir, "nature", "init.lua"))
		if err2 != nil {
			fmt.Fprintln(os.Stderr, "Missing nature module, some functionality and builtins will be missing.")
			fmt.Fprintln(os.Stderr, "local error:", err1)
			fmt.Fprintln(os.Stderr, "global install error:", err2)
		}
	}
}

func runConfig(confpath string) {
	if !interactive {
		return
	}
	err := l.DoFile(confpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "\nAn error has occured while loading your config! Falling back to minimal default config.")
		l.DoString(minimalconf)
	}
}
