// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
	"os"
	"path/filepath"
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
	"read": hlread,
}

func hilbishLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)

	host, _ := os.Hostname()
	username := curuser.Username
	if runtime.GOOS == "windows" {
		username = strings.Split(username, "\\")[1] // for some reason Username includes the hostname on windows
	}

	util.SetField(L, mod, "ver", lua.LString(version), "Hilbish version")
	util.SetField(L, mod, "user", lua.LString(username), "Username of user")
	util.SetField(L, mod, "host", lua.LString(host), "Host name of the machine")
	util.SetField(L, mod, "home", lua.LString(curuser.HomeDir), "Home directory of the user")
	util.SetField(L, mod, "dataDir", lua.LString(dataDir), "Directory for Hilbish's data files")
	util.SetField(L, mod, "interactive", lua.LBool(interactive), "If this is an interactive shell")
	util.SetField(L, mod, "login", lua.LBool(interactive), "Whether this is a login shell")
	util.Document(L, mod, "Hilbish's core API, containing submodules and functions which relate to the shell itself.")

	// hilbish.userDir table
	hshuser := L.NewTable()
	userConfigDir, _ := os.UserConfigDir()
	userDataDir := ""
	// i honestly dont know what directories to use for this
	switch runtime.GOOS {
	case "linux":
		userDataDir = getenv("XDG_DATA_HOME", curuser.HomeDir + "/.local/share")
	default:
		userDataDir = filepath.Join(userConfigDir)
	}

	util.SetField(L, hshuser, "config", lua.LString(userConfigDir), "User's config directory")
	util.SetField(L, hshuser, "data", lua.LString(userDataDir), "XDG data directory")
	util.Document(L, hshuser, "User directories to store configs and/or modules.")
	L.SetField(mod, "userDir", hshuser)

	// hilbish.aliases table
	aliases = NewAliases()
	aliasesModule := aliases.Loader(L)
	util.Document(L, aliasesModule, "Alias inferface for Hilbish.")
	L.SetField(mod, "aliases", aliasesModule)

	L.Push(mod)

	return 1
}

// run(cmd)
// Runs `cmd` in Hilbish's sh interpreter.
func hlrun(L *lua.LState) int {
	var exitcode uint8
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

	flag := getopt.Lookup([]rune(flagchar)[0])
	if flag == nil {
		L.Push(lua.LNil)
		return 1
	}

	passed := flag.Seen()
	L.Push(lua.LBool(passed))

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

// read(prompt) -> input?
// Read input from the user, using Hilbish's line editor/input reader.
// This is a separate instance from the one Hilbish actually uses.
// Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)
func hlread(L *lua.LState) int {
	luaprompt := L.CheckString(1)
	lualr := newLineReader(luaprompt)

	input, err := lualr.Read()
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(input))
	return 1
}

