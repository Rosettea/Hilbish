// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
	"os"
	"runtime"
	"strings"

	"hilbish/util"

	"github.com/pborman/getopt"
	"github.com/yuin/gopher-lua"
	"mvdan.cc/sh/v3/interp"
)

var exports = map[string]lua.LGFunction {
	"run": hlrun,
	"flag": hlflag,
	"cwd": hlcwd,
}

func HilbishLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)

	host, _ := os.Hostname()
	username := curuser.Username
	// this will be baked into binary since GOOS is a constant
	if runtime.GOOS == "windows" {
		username = strings.Split(username, "\\")[1] // for some reason Username includes the hostname on windows
	}

	L.SetField(mod, "ver", lua.LString(version))
	L.SetField(mod, "user", lua.LString(username))
	L.SetField(mod, "host", lua.LString(host))
	L.SetField(mod, "home", lua.LString(homedir))

	xdg := L.NewTable()
	L.SetField(xdg, "config", lua.LString(confDir))
	L.SetField(xdg, "data", lua.LString(getenv("XDG_DATA_HOME", homedir + "/.local/share/")))
	L.SetField(mod, "xdg", xdg)

	util.Document(L, mod, "A miscellaneous sort of \"core\" API for things that relate to the shell itself and others.")
	L.Push(mod)

	return 1
}

// run(cmd)
// Runs `cmd` in Hilbish's sh interpreter
func hlrun(L *lua.LState) int {
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

// flag(f)
// Checks if the `f` flag has been passed to Hilbish.
func hlflag(L *lua.LState) int {
	flagchar := L.CheckString(1)

	L.Push(lua.LBool(getopt.Lookup([]rune(flagchar)[0]).Seen()))

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
