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

var rtmm *rt.Runtime
var watcherMetaKey = rt.StringValue("hshwatcher")
var Loader = packagelib.Loader{
	Load: loaderFunc,
	Name: "fs",
}

func loaderFunc(rtm *rt.Runtime) (rt.Value, func()) {
	rtmm = rtm
	watcherMethods := rt.NewTable()
	watcherFuncs := map[string]util.LuaExport{
		"start": {watcherStart, 1, false},
		"stop": {watcherStop, 1, false},
	}
	util.SetExports(rtm, watcherMethods, watcherFuncs)

	watcherMeta := rt.NewTable()
	watcherIndex := func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		//ti, _ := watcherArg(c, 0)

		arg := c.Arg(1)
		val := watcherMethods.Get(arg)

		if val != rt.NilValue {
			return c.PushingNext1(t.Runtime, val), nil
		}

		return c.PushingNext1(t.Runtime, val), nil
	}

	watcherMeta.Set(rt.StringValue("__index"), rt.FunctionValue(rt.NewGoFunction(watcherIndex, "__index", 2, false)))
	rtm.SetRegistry(watcherMetaKey, rt.TableValue(watcherMeta))

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
		"watch": util.LuaExport{fwatch, 2, false},
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

// stat(path) -> {}
// Returns a table of info about the `path`.
// It contains the following keys:
// name (string) - Name of the path
// size (number) - Size of the path
// mode (string) - Permission mode in an octal format string (with leading 0)
// isDir (boolean) - If the path is a directory
// --- @param path string
// --- @returns table
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

// readdir(dir) -> {}
// Returns a table of files in `dir`.
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

// abs(path) -> string
// Gives an absolute version of `path`.
// --- @param path string
// --- @returns string
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

// basename(path) -> string
// Gives the basename of `path`. For the rules,
// see Go's filepath.Base
// --- @returns string
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

// dir(path) -> string
// Returns the directory part of `path`. For the rules, see Go's
// filepath.Dir
// --- @param path string
// --- @returns string
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

// glob(pattern) -> matches (table)
// Glob all files and directories that match the pattern.
// For the rules, see Go's filepath.Glob
// --- @param pattern string
// --- @returns table
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

// join(...) -> string
// Takes paths and joins them together with the OS's
// directory separator (forward or backward slash).
// --- @vararg string
// --- @returns string
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

func fwatch(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}

	dir, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	watcher, err := c.ClosureArg(1)
	if err != nil {
		return nil, err
	}

	dw := newWatcher(dir, watcher)

	return c.PushingNext1(t.Runtime, rt.UserDataValue(dw.ud)), nil
}
