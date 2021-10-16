// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
	"os"
	"runtime"
	"strings"

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

	setField(L, mod, "ver", lua.LString(version), "The version of Hilbish")
	setField(L, mod, "user", lua.LString(username), "Current user's username")
	setField(L, mod, "host", lua.LString(host), "Hostname of the system")
	setField(L, mod, "home", lua.LString(homedir), "Path of home directory")
	setField(L, mod, "dataDir", lua.LString(dataDir), "Path of Hilbish's data files")

	xdg := L.NewTable()
	setField(L, xdg, "config", lua.LString(confDir), "XDG config directory")
	setField(L, xdg, "data", lua.LString(getenv("XDG_DATA_HOME", homedir + "/.local/share/")), "XDG data directory")
	setField(L, mod, "xdg", xdg, "XDG values for Linux")

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
