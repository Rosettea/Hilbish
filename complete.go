package main

import (
	"errors"
	"path/filepath"
	"strings"
	"os"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

var completer rt.Value

func fileComplete(query, ctx string, fields []string) ([]string, string) {
	return matchPath(query)
}

func binaryComplete(query, ctx string, fields []string) ([]string, string) {
	var completions []string

	prefixes := []string{"./", "../", "/", "~/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(query, prefix) {
			fileCompletions, filePref := matchPath(query)
			if len(fileCompletions) != 0 {
				for _, f := range fileCompletions {
					fullPath, _ := filepath.Abs(expandHome(query + strings.TrimPrefix(f, filePref)))
					if err := findExecutable(fullPath, false, true); err != nil {
						continue
					}
					completions = append(completions, f)
				}
			}
			return completions, filePref
		}
	}

	// filter out executables, but in path
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		// print dir to stderr for debugging
		// search for an executable which matches our query string
		if matches, err := filepath.Glob(filepath.Join(dir, query + "*")); err == nil {
			// get basename from matches
			for _, match := range matches {
				// check if we have execute permissions for our match
				err := findExecutable(match, true, false)
				if err != nil {
					continue
				}
				// get basename from match
				name := filepath.Base(match)
				// add basename to completions
				completions = append(completions, name)
			}
		}
	}

	// add lua registered commands to completions
	for cmdName := range commands {
		if strings.HasPrefix(cmdName, query) {
			completions = append(completions, cmdName)
		}
	}

	completions = removeDupes(completions)

	return completions, query
}

func matchPath(query string) ([]string, string) {
	var entries []string
	var baseName string

	path, _ := filepath.Abs(expandHome(filepath.Dir(query)))
	if string(query) == "" {
		// filepath base below would give us "."
		// which would cause a match of only dotfiles
		path, _ = filepath.Abs(".")
	} else if !strings.HasSuffix(query, string(os.PathSeparator)) {
		baseName = filepath.Base(query)
	}

	files, _ := os.ReadDir(path)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), baseName) {
			entry := file.Name()
			if file.IsDir() {
				entry = entry + string(os.PathSeparator)
			}
			entry = escapeFilename(entry)
			entries = append(entries, entry)
		}
	}

	return entries, baseName
}

func escapeFilename(fname string) string {
	args := []string{
		"\"", "\\\"",
		"'", "\\'",
		"`", "\\`",
		" ", "\\ ",
		"(", "\\(",
		")", "\\)",
		"[", "\\[",
		"]", "\\]",
		"$", "\\$",
		"&", "\\&",
		"*", "\\*",
		">", "\\>",
		"<", "\\<",
		"|", "\\|",
	}

	r := strings.NewReplacer(args...)
	return r.Replace(fname)
}

func completionLoader(rtm *rt.Runtime) *rt.Table {
	exports := map[string]util.LuaExport{
		"files": {luaFileComplete, 3, false},
		"bins": {luaBinaryComplete, 3, false},
		"call": {callLuaCompleter, 4, false},
		"handler": {completionHandler, 2, false},
	}

	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)
	
	return mod
}

func completionHandler(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	line, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	// just for validation
	_, err = c.IntArg(1)
	if err != nil {
		return nil, err
	}

	ctx := strings.TrimLeft(line, " ")
	if len(ctx) == 0 {
		return c.PushingNext(t.Runtime, rt.TableValue(rt.NewTable()), rt.StringValue("")), nil
	}

	ctx = aliases.Resolve(ctx)
	fields := strings.Split(ctx, " ")
	query := fields[len(fields) - 1]

	luaFields := rt.NewTable()

	for i, f := range fields {
		luaFields.Set(rt.IntValue(int64(i + 1)), rt.StringValue(f))
	}

	compMod := hshMod.Get(rt.StringValue("completion")).AsTable()
	var term *rt.Termination
	if len(fields) == 1 {
		term = rt.NewTerminationWith(t.CurrentCont(), 2, false)
		err := rt.Call(t, compMod.Get(rt.StringValue("bins")), []rt.Value{
			rt.StringValue(query),
			rt.StringValue(ctx),
			rt.TableValue(luaFields),
		}, term)

		if err != nil {
			return nil, err
		}
	} else {
		gterm := rt.NewTerminationWith(t.CurrentCont(), 2, false)
		err := rt.Call(t, compMod.Get(rt.StringValue("call")), []rt.Value{
			rt.StringValue("commands." + fields[0]),
			rt.StringValue(query),
			rt.StringValue(ctx),
			rt.TableValue(luaFields),
		}, gterm)

		if err == nil {
			groups := gterm.Get(0)
			pfx := gterm.Get(1)

			return c.PushingNext(t.Runtime, groups, pfx), nil
		}

		// error means there isnt a command handler - default to files in that case
		term = rt.NewTerminationWith(t.CurrentCont(), 2, false)
		err = rt.Call(t, compMod.Get(rt.StringValue("files")), []rt.Value{
			rt.StringValue(query),
			rt.StringValue(ctx),
			rt.TableValue(luaFields),
		}, term)
	}

	comps := term.Get(0)
	pfx := term.Get(1)

	groups := rt.NewTable()

	compGroup := rt.NewTable()
	compGroup.Set(rt.StringValue("items"), comps)
	compGroup.Set(rt.StringValue("type"), rt.StringValue("grid"))

	groups.Set(rt.IntValue(1), rt.TableValue(compGroup))
	return c.PushingNext(t.Runtime, rt.TableValue(groups), pfx), nil
}

func callLuaCompleter(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(4); err != nil {
		return nil, err
	}
	completer, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	query, err := c.StringArg(1)
	if err != nil {
		return nil, err
	}
	ctx, err := c.StringArg(2)
	if err != nil {
		return nil, err
	}
	fields, err := c.TableArg(3)
	if err != nil {
		return nil, err
	}

	var completecb *rt.Closure
	var ok bool
	if completecb, ok = luaCompletions[completer]; !ok {
		return nil, errors.New("completer " + completer + " does not exist")
	}

	// we must keep the holy 80 cols
	completerReturn, err := rt.Call1(l.MainThread(),
	rt.FunctionValue(completecb), rt.StringValue(query),
	rt.StringValue(ctx), rt.TableValue(fields))

	if err != nil {
		return nil, err
	}

	return c.PushingNext1(t.Runtime, completerReturn), nil
}

func luaFileComplete(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	query, ctx, fds, err := getCompleteParams(t, c)
	if err != nil {
		return nil, err
	}

	completions, pfx := fileComplete(query, ctx, fds)
	luaComps := rt.NewTable()

	for i, comp := range completions {
		luaComps.Set(rt.IntValue(int64(i + 1)), rt.StringValue(comp))
	}

	return c.PushingNext(t.Runtime, rt.TableValue(luaComps), rt.StringValue(pfx)), nil
}

func luaBinaryComplete(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	query, ctx, fds, err := getCompleteParams(t, c)
	if err != nil {
		return nil, err
	}

	completions, pfx := binaryComplete(query, ctx, fds)
	luaComps := rt.NewTable()

	for i, comp := range completions {
		luaComps.Set(rt.IntValue(int64(i + 1)), rt.StringValue(comp))
	}

	return c.PushingNext(t.Runtime, rt.TableValue(luaComps), rt.StringValue(pfx)), nil
}

func getCompleteParams(t *rt.Thread, c *rt.GoCont) (string, string, []string, error) {
	if err := c.CheckNArgs(3); err != nil {
		return "", "", []string{}, err
	}
	query, err := c.StringArg(0)
	if err != nil {
		return "", "", []string{}, err
	}
	ctx, err := c.StringArg(1)
	if err != nil {
		return "", "", []string{}, err
	}
	fields, err := c.TableArg(2)
	if err != nil {
		return "", "", []string{}, err
	}

	var fds []string
	util.ForEach(fields, func(k rt.Value, v rt.Value) {
		if v.Type() == rt.StringType {
			fds = append(fds, v.AsString())
		}
	})

	return query, ctx, fds, err
}
