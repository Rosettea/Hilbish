// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
	"os"

	"github.com/pborman/getopt"
	"github.com/yuin/gopher-lua"
	"mvdan.cc/sh/v3/interp"
)

var exports = map[string]lua.LGFunction {
	"run": run,
	"flag": flag,
	"cwd": cwd,
}

func HilbishLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)

	host, _ := os.Hostname()

	L.SetField(mod, "ver", lua.LString(version))
	L.SetField(mod, "user", lua.LString(curuser.Username))
	L.SetField(mod, "host", lua.LString(host))

	L.Push(mod)

	return 1
}

// Runs a command
func run(L *lua.LState) int {
	var exitcode uint8 = 0
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

func flag(L *lua.LState) int {
	flagchar := L.CheckString(1)

	L.Push(lua.LBool(getopt.Lookup([]rune(flagchar)[0]).Seen()))

	return 1
}

func cwd(L *lua.LState) int {
	cwd, _ := os.Getwd()

	L.Push(lua.LString(cwd))

	return 1
}

