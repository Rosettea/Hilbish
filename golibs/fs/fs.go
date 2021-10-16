// The fs module provides easy and simple access to filesystem functions and other
// things, and acts an addition to the Lua standard library's I/O and fs functions.
package fs

import (
	"fmt"
	"os"
	"strings"

	"hilbish/util"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)

	util.Document(L, mod, `The fs module provides easy and simple access to filesystem functions and other
things, and acts an addition to the Lua standard library's I/O and fs functions.`)

	L.Push(mod)
	return 1
}

func luaErr(L *lua.LState, msg string) {
	L.Error(lua.LString(msg), 2)
}

var exports = map[string]lua.LGFunction{
	"cd": fcd,
	"mkdir": fmkdir,
	"stat": fstat,
}

// cd(dir)
// Changes directory to `dir`
func fcd(L *lua.LState) int {
	path := L.CheckString(1)

	err := os.Chdir(strings.TrimSpace(path))
	if err != nil {
		switch e := err.(*os.PathError).Err.Error(); e {
		case "no such file or directory":
			luaErr(L, 1)
		case "not a directory":
			luaErr(L, 2)
		default:
			fmt.Printf("Found unhandled error case: %s\n", e)
			fmt.Printf("Report this at https://github.com/Rosettea/Hilbish/issues with the title being: \"fs: unhandled error case %s\", and show what caused it.\n", e)
			luaErr(L, 213)
		}
	}

	return 0
}

// mkdir(name, recursive)
// Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
func fmkdir(L *lua.LState) int {
	dirname := L.CheckString(1)
	recursive := L.ToBool(2)
	path := strings.TrimSpace(dirname)

	// TODO: handle error here
	if recursive {
		os.MkdirAll(path, 0744)
	} else {
		os.Mkdir(path, 0744)
	}

	return 0
}

// stat(path)
// Returns info about `path`
func fstat(L *lua.LState) int {
	path := L.CheckString(1)

	// TODO: handle error here
	pathinfo, _ := os.Stat(path)
	L.Push(luar.New(L, pathinfo))

	return 1
}
