package main

import (
	"errors"
	"path/filepath"
	"strings"
	"os"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

var charEscapeMap = []string{
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
var charEscapeMapInvert = invert(charEscapeMap)
var escapeReplaer = strings.NewReplacer(charEscapeMap...)
var escapeInvertReplaer = strings.NewReplacer(charEscapeMapInvert...)

func invert(m []string) []string {
	newM := make([]string, len(charEscapeMap))
	for i := range m {
		if (i + 1) % 2 == 0 {
			newM[i] = m[i - 1]
			newM[i - 1] = m[i]
		}
	}

	return newM
}

func splitForFile(str string) []string {
	split := []string{}
	sb := &strings.Builder{}
	quoted := false

	for i, r := range str {
		if r == '"' {
			quoted = !quoted
			sb.WriteRune(r)
		} else if r == ' ' && str[i - 1] == '\\' {
			sb.WriteRune(r)
		} else if !quoted && r == ' ' {
			split = append(split, sb.String())
			sb.Reset()
		} else {
			sb.WriteRune(r)
		}
	}
	if strings.HasSuffix(str, " ") {
		split = append(split, "")
	}

	if sb.Len() > 0 {
		split = append(split, sb.String())
	}

	return split
}

func fileComplete(query, ctx string, fields []string) ([]string, string) {
	q := splitForFile(ctx)
	path := ""
	if len(q) != 0 {
		path = q[len(q) - 1]
	}

	return matchPath(path)
}

func binaryComplete(query, ctx string, fields []string) ([]string, string) {
	q := splitForFile(ctx)
	query = ""
	if len(q) != 0 {
		query = q[len(q) - 1]
	}

	var completions []string

	prefixes := []string{"./", "../", "/", "~/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(query, prefix) {
			fileCompletions, filePref := matchPath(query)
			if len(fileCompletions) != 0 {
				for _, f := range fileCompletions {
					fullPath, _ := filepath.Abs(util.ExpandHome(query + strings.TrimPrefix(f, filePref)))
					if err := findExecutable(escapeInvertReplaer.Replace(fullPath), false, true); err != nil {
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
	oldQuery := query
	query = strings.TrimPrefix(query, "\"")
	var entries []string
	var baseName string

	query = escapeInvertReplaer.Replace(query)
	path, _ := filepath.Abs(util.ExpandHome(filepath.Dir(query)))
	if string(query) == "" {
		// filepath base below would give us "."
		// which would cause a match of only dotfiles
		path, _ = filepath.Abs(".")
	} else if !strings.HasSuffix(query, string(os.PathSeparator)) {
		baseName = filepath.Base(query)
	}

	files, _ := os.ReadDir(path)
	for _, entry := range files {
		// should we handle errors here?
		file, err := entry.Info()
		if err == nil && file.Mode() & os.ModeSymlink != 0 {
			path, err := filepath.EvalSymlinks(filepath.Join(path, file.Name()))
			if err == nil {
				file, err = os.Lstat(path)
			}
		}

		if strings.HasPrefix(file.Name(), baseName) {
			entry := file.Name()
			if file.IsDir() {
				entry = entry + string(os.PathSeparator)
			}
			if !strings.HasPrefix(oldQuery, "\"") {
				entry = escapeFilename(entry)
			}
			entries = append(entries, entry)
		}
	}
	if !strings.HasPrefix(oldQuery, "\"") {
		baseName = escapeFilename(baseName)
	}

	return entries, baseName
}

func escapeFilename(fname string) string {
	return escapeReplaer.Replace(fname)
}

// #interface completions
// tab completions
// The completions interface deals with tab completions.
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

// #interface completions
// handler(line, pos)
// The handler function is the callback for tab completion in Hilbish.
// You can check the completions doc for more info.
// --- @param line string
// --- @param pos string
func completionHandler(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	return c.Next(), nil
}

// #interface completions
// call(name, query, ctx, fields) -> completionGroups (table), prefix (string)
// Calls a completer function. This is mainly used to call
// a command completer, which will have a `name` in the form
// of `command.name`, example: `command.git`.
// You can check `doc completions` for info on the `completionGroups` return value.
// --- @param name string
// --- @param query string
// --- @param ctx string
// --- @param fields table
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

// #interface completions
// files(query, ctx, fields) -> entries (table), prefix (string)
// Returns file completion candidates based on the provided query.
// --- @param query string
// --- @param ctx string
// --- @param fields table
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

// #interface completions
// bins(query, ctx, fields) -> entries (table), prefix (string)
// Returns binary/executale completion candidates based on the provided query.
// --- @param query string
// --- @param ctx string
// --- @param fields table
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
