package fs

import (
	"os"

	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
    mod := L.SetFuncs(L.NewTable(), exports)

    L.Push(mod)
    return 1
}


func LuaErr(L *lua.LState, code int) {
	// TODO: Error with a table, with path and error code
	L.Error(lua.LNumber(code), 1)
}

var exports = map[string]lua.LGFunction{
    "cd": cd,
}

func cd(L *lua.LState) int {
	path := L.ToString(1)

	err := os.Chdir(path)
	if err != nil {
		switch err.(*os.PathError).Err.Error() {
		case "no such file or directory":
			LuaErr(L, 2)
		}
	}

	return 0
}

