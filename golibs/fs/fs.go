// The fs module provides easy and simple access to filesystem functions and other
// things, and acts an addition to the Lua standard library's I/O and fs functions.
package fs

import (
	"strconv"
	"os"
	"strings"

	"hilbish/util"
	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)

	util.Document(L, mod, `The fs module provides easy and simple access to
filesystem functions and other things, and acts an
addition to the Lua standard library's I/O and fs functions.`)

	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"cd": fcd,
	"mkdir": fmkdir,
	"stat": fstat,
	"readdir": freaddir,
}

// cd(dir)
// Changes directory to `dir`
// --- @param dir string
func fcd(L *lua.LState) int {
	path := L.CheckString(1)

	err := os.Chdir(strings.TrimSpace(path))
	if err != nil {
		e := err.(*os.PathError).Err.Error()
		L.RaiseError(e + ": " + path)
	}

	return 0
}

// mkdir(name, recursive)
// Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
// --- @param name string
// --- @param recursive boolean
func fmkdir(L *lua.LState) int {
	dirname := L.CheckString(1)
	recursive := L.ToBool(2)
	path := strings.TrimSpace(dirname)
	var err error

	if recursive {
		err = os.MkdirAll(path, 0744)
	} else {
		err = os.Mkdir(path, 0744)
	}
	if err != nil {
		L.RaiseError(err.Error() + ": " + path)
	}

	return 0
}

// stat(path)
// Returns info about `path`
// --- @param path string
func fstat(L *lua.LState) int {
	path := L.CheckString(1)

	pathinfo, err := os.Stat(path)
	if err != nil {
		L.RaiseError(err.Error() + ": " + path)
		return 0
	}
	statTbl := L.NewTable()
	L.SetField(statTbl, "name", lua.LString(pathinfo.Name()))
	L.SetField(statTbl, "size", lua.LNumber(pathinfo.Size()))
	L.SetField(statTbl, "mode", lua.LString("0" + strconv.FormatInt(int64(pathinfo.Mode().Perm()), 8)))
	L.SetField(statTbl, "isDir", lua.LBool(pathinfo.IsDir()))
	L.Push(statTbl)

	return 1
}

// readdir(dir)
// Returns a table of files in `dir`
// --- @param dir string
// --- @return table
func freaddir(L *lua.LState) int {
	dir := L.CheckString(1)
	names := L.NewTable()

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		L.RaiseError(err.Error() + ": " + dir)
		return 0
	}
	for _, entry := range dirEntries {
		names.Append(lua.LString(entry.Name()))
	}

	L.Push(names)

	return 1
}
