package fs

import (
	"strconv"
	"os"
	"strings"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib/packagelib"
)

var Loader = packagelib.Loader{
	Load: loaderFunc,
	Name: "fs",
}

func loaderFunc(rtm *rt.Runtime) (rt.Value, func()) {
	exports := map[string]util.LuaExport{
		"cd": util.LuaExport{fcd, 1, false},
		"mkdir": util.LuaExport{fmkdir, 2, false},
		"stat": util.LuaExport{fstat, 1, false},
		"readdir": util.LuaExport{freaddir, 1, false},
	}
	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	util.Document(mod, `The fs module provides easy and simple access to
filesystem functions and other things, and acts an
addition to the Lua standard library's I/O and fs functions.`)

	return rt.TableValue(mod), nil
}

// cd(dir)
// Changes directory to `dir`
// --- @param dir string
func fcd(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	err = os.Chdir(strings.TrimSpace(path))
	if err != nil {
		return nil, err
	}

	return c.Next(), err
}

// mkdir(name, recursive)
// Makes a directory called `name`. If `recursive` is true, it will create its parent directories.
// --- @param name string
// --- @param recursive boolean
func fmkdir(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	dirname, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	recursive, err := c.BoolArg(1)
	if err != nil {
		return nil, err
	}
	path := strings.TrimSpace(dirname)

	if recursive {
		err = os.MkdirAll(path, 0744)
	} else {
		err = os.Mkdir(path, 0744)
	}
	if err != nil {
		return nil, err
	}

	return c.Next(), err
}

// stat(path)
// Returns info about `path`
// --- @param path string
func fstat(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	pathinfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	statTbl := rt.NewTable()
	statTbl.Set(rt.StringValue("name"), rt.StringValue(pathinfo.Name()))
	statTbl.Set(rt.StringValue("size"), rt.IntValue(pathinfo.Size()))
	statTbl.Set(rt.StringValue("mode"), rt.StringValue("0" + strconv.FormatInt(int64(pathinfo.Mode().Perm()), 8)))
	statTbl.Set(rt.StringValue("isDir"), rt.BoolValue(pathinfo.IsDir()))
	
	return c.PushingNext1(t.Runtime, rt.TableValue(statTbl)), nil
}

// readdir(dir)
// Returns a table of files in `dir`
// --- @param dir string
// --- @return table
func freaddir(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	dir, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	names := rt.NewTable()

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for i, entry := range dirEntries {
		names.Set(rt.IntValue(int64(i + 1)), rt.StringValue(entry.Name()))
	}

	return c.PushingNext1(t.Runtime, rt.TableValue(names)), nil
}
