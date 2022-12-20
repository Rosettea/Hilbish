// filesystem interaction and functionality library
// The fs module provides easy and simple access to filesystem functions
// and other things, and acts an addition to the Lua standard library's
// I/O and filesystem functions.
package fs

import (
	"fmt"
	"path/filepath"
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
		"abs": util.LuaExport{fabs, 1, false},
		"basename": util.LuaExport{fbasename, 1, false},
		"dir": util.LuaExport{fdir, 1, false},
		"glob": util.LuaExport{fglob, 1, false},
		"join": util.LuaExport{fjoin, 0, true},
	}
	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)
	mod.Set(rt.StringValue("pathSep"), rt.StringValue(string(os.PathSeparator)))
	mod.Set(rt.StringValue("pathListSep"), rt.StringValue(string(os.PathListSeparator)))

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
	path = util.ExpandHome(strings.TrimSpace(path))

	err = os.Chdir(path)
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
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	recursive, err := c.BoolArg(1)
	if err != nil {
		return nil, err
	}
	path = util.ExpandHome(strings.TrimSpace(path))

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
	path = util.ExpandHome(path)

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
	dir = util.ExpandHome(dir)
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

// abs(path)
// Gives an absolute version of `path`.
// --- @param path string
func fabs(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	path = util.ExpandHome(path)

	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return c.PushingNext1(t.Runtime, rt.StringValue(abspath)), nil
}

// basename(path)
// Gives the basename of `path`. For the rules,
// see Go's filepath.Base
func fbasename(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	return c.PushingNext(t.Runtime, rt.StringValue(filepath.Base(path))), nil
}

// dir(path)
// Returns the directory part of `path`. For the rules, see Go's
// filepath.Dir
func fdir(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	return c.PushingNext(t.Runtime, rt.StringValue(filepath.Dir(path))), nil
}

// glob(pattern)
// Glob all files and directories that match the pattern.
// For the rules, see Go's filepath.Glob
func fglob(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	pattern, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	luaMatches := rt.NewTable()

	for i, match := range matches {
		luaMatches.Set(rt.IntValue(int64(i + 1)), rt.StringValue(match))
	}
	
	return c.PushingNext(t.Runtime, rt.TableValue(luaMatches)), nil
}

// join(paths...)
// Takes paths and joins them together with the OS's
// directory separator (forward or backward slash).
func fjoin(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	strs := make([]string, len(c.Etc()))
	for i, v := range c.Etc() {
		if v.Type() != rt.StringType {
			// +2; go indexes of 0 and first arg from above
			return nil, fmt.Errorf("bad argument #%d to run (expected string, got %s)", i + 1, v.TypeName())
		}
		strs[i] = v.AsString()
	}

	res := filepath.Join(strs...)

	return c.PushingNext(t.Runtime, rt.StringValue(res)), nil
}
