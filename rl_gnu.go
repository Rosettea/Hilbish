// +build gnurl

package main

// Here we define a generic interface for readline and hilbiline,
// making them interchangable during build time
// this is normal readline

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Rosettea/readline"
	"github.com/yuin/gopher-lua"
)

type lineReader struct {
	Prompt string
}

func newLineReader(prompt string) *lineReader {
	readline.Init()

	readline.Completer = func(query string, ctx string) []string {
		var completions []string
		// trim whitespace from ctx
		ctx = strings.TrimLeft(ctx, " ")
		if len(ctx) == 0 {
			return []string{}
		}
		fields := strings.Split(ctx, " ")
		if len(fields) == 0 {
			return []string{}
		}

		ctx = aliases.Resolve(ctx)

		if len(fields) == 1 {
			fileCompletions := fileComplete(query, ctx, fields)
			if len(fileCompletions) != 0 {
				return fileCompletions
			}

			// filter out executables, but in path
			for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
				// print dir to stderr for debugging
				// search for an executable which matches our query string
				if matches, err := filepath.Glob(filepath.Join(dir, query + "*")); err == nil {
					// get basename from matches
					for _, match := range matches {
						// check if we have execute permissions for our match
						if info, err := os.Stat(match); err == nil && info.Mode().Perm() & 0100 == 0 {
							continue
						}
						// get basename from match
						name := filepath.Base(match)
						// print name to stderr for debugging
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
		} else {
			if completecb, ok := luaCompletions["command." + fields[0]]; ok {
				err := l.CallByParam(lua.P{
					Fn: completecb,
					NRet: 1,
					Protect: true,
				})

				if err != nil {
					return []string{}
				}

				luacompleteTable := l.Get(-1)
				l.Pop(1)

				if cmpTbl, ok := luacompleteTable.(*lua.LTable); ok {
					cmpTbl.ForEach(func(key lua.LValue, value lua.LValue) {
						// if key is a number (index), we just check and complete that
						if key.Type() == lua.LTNumber {
							// if we have only 2 fields then this is fine
							if len(fields) == 2 {
								if strings.HasPrefix(value.String(), fields[1]) {
									completions = append(completions, value.String())
								}
							}
						} else if key.Type() == lua.LTString {
							if len(fields) == 2 {
								if strings.HasPrefix(key.String(), fields[1]) {
									completions = append(completions, key.String())
								}
							} else {
								// if we have more than 2 fields, we need to check if the key matches
								// the current field and if it does, we need to check if the value is a string
								// or table (nested sub completions)
								if key.String() == fields[1] {
									// if value is a table, we need to iterate over it
									// and add each value to completions
									// check if value is either a table or function
									if value.Type() == lua.LTTable {
										valueTbl := value.(*lua.LTable)
										valueTbl.ForEach(func(key lua.LValue, value lua.LValue) {
											val := value.String()
											if val == "<file>" {
												// complete files
												completions = append(completions, readline.FilenameCompleter(query, ctx)...)
											} else {
												if strings.HasPrefix(val, query) {
													completions = append(completions, val)
												}
											}
										})
									} else if value.Type() == lua.LTFunction {
										// if value is a function, we need to call it
										// and add each value to completions
										// completionsCtx is the context we pass to the function,
										// removing 2 fields from the fields array
										completionsCtx := strings.Join(fields[2:], " ")
										err := l.CallByParam(lua.P{
											Fn: value,
											NRet: 1,
											Protect: true,
										}, lua.LString(query), lua.LString(completionsCtx))

										if err != nil {
											return
										}

										luacompleteTable := l.Get(-1)
										l.Pop(1)

										// just check if its actually a table and add it to the completions
										if cmpTbl, ok := luacompleteTable.(*lua.LTable); ok {
											cmpTbl.ForEach(func(key lua.LValue, value lua.LValue) {
												val := value.String()
												if strings.HasPrefix(val, query) {
													completions = append(completions, val)
												}
											})
										}
									} else {
										// throw lua error
										// complete.cmdname: error message...
										l.RaiseError("complete." + fields[0] + ": completion value is not a table or function")
									}
								}
							}
						}
					})
				}
			}

			if len(completions) == 0 {
				completions = readline.FilenameCompleter(query, ctx)
			}
		}
		return completions
	}
	readline.LoadHistory(defaultHistPath)

	return &lineReader{
		Prompt: prompt,
	}
}

func (lr *lineReader) Read() (string, error) {
	hooks.Em.Emit("command.precmd", nil)
	return readline.String(lr.Prompt)
}

func (lr *lineReader) SetPrompt(prompt string) {
	lr.Prompt = prompt
}

func (lr *lineReader) AddHistory(cmd string) {
	readline.AddHistory(cmd)
	readline.SaveHistory(defaultHistPath)
}

func (lr *lineReader) ClearInput() {
	readline.ReplaceLine("", 0)
	readline.RefreshLine()
}

func (lr *lineReader) Resize() {
	readline.Resize()
}

// lua module
func (lr *lineReader) Loader(L *lua.LState) *lua.LTable {
	lrLua := map[string]lua.LGFunction{
		"add": lr.luaAddHistory,
		"all": lr.luaAllHistory,
		"clear": lr.luaClearHistory,
		"get": lr.luaGetHistory,
		"size": lr.luaSize,
	}

	mod := L.SetFuncs(L.NewTable(), lrLua)

	return mod
}

func (lr *lineReader) luaAddHistory(l *lua.LState) int {
	cmd := l.CheckString(1)
	lr.AddHistory(cmd)

	return 0
}

func (lr *lineReader) luaSize(l *lua.LState) int {
	l.Push(lua.LNumber(readline.HistorySize()))

	return 1
}

func (lr *lineReader) luaGetHistory(l *lua.LState) int {
	idx := l.CheckInt(1)
	cmd := readline.GetHistory(idx)
	l.Push(lua.LString(cmd))

	return 1
}

func (lr *lineReader) luaAllHistory(l *lua.LState) int {
	tbl := l.NewTable()
	size := readline.HistorySize()

	for i := 0; i < size; i++ {
		cmd := readline.GetHistory(i)
		tbl.Append(lua.LString(cmd))
	}

	l.Push(tbl)

	return 1
}

func (lr *lineReader) luaClearHistory(l *lua.LState) int {
	readline.ClearHistory()
	readline.SaveHistory(defaultHistPath)

	return 0
}
