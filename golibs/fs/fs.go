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

var exports = map[string]lua.LGFunction{
    "cd": cd,
}

func cd(L *lua.LState) int {
	path := L.ToString(1)

	os.Chdir(path)

	return 0
}

