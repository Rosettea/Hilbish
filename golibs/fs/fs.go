package fs

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}

func LuaErr(L *lua.LState, code int) {
	// TODO: Error with a table, with path and error code
	L.Error(lua.LNumber(code), 2)
}

var exports = map[string]lua.LGFunction{
	"cd": cd,
	"mkdir": mkdir,
	"stat": stat,
}

func cd(L *lua.LState) int {
	path := L.CheckString(1)

	err := os.Chdir(strings.TrimSpace(path))
	if err != nil {
		switch e := err.(*os.PathError).Err.Error(); e {
		case "no such file or directory":
			LuaErr(L, 1)
		default:
			fmt.Printf("Found unhandled error case: %s", e)
			fmt.Printf("Report this at https://github.com/Hilbis/Hilbish/issues with the title being: \"fs: unahndled error case %s\", and show what caused it.\n", e)
			LuaErr(L, 213)
		}
	}

	return 0
}

func mkdir(L *lua.LState) int {
	dirname := L.CheckString(1)

	// TODO: handle error here
	os.Mkdir(strings.TrimSpace(dirname), 0744)

	return 0
}

func stat(L *lua.LState) int {
	path := L.CheckString(1)

	// TODO: handle error here
	pathinfo, _ := os.Stat(path)
	L.Push(luar.New(L, pathinfo))

	return 1
}
