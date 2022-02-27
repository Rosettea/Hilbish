// +build !gnurl

package main

// Here we define a generic interface for readline and hilbiline,
// making them interchangable during build time
// this is hilbiline's, as is obvious by the filename

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"os"

	"github.com/maxlandon/readline"
	"github.com/yuin/gopher-lua"
)

type lineReader struct {
	rl *readline.Instance
}

// other gophers might hate this naming but this is local, shut up
func newLineReader(prompt string) *lineReader {
	rl := readline.NewInstance()
	rl.Multiline = true
	rl.TabCompleter = func(line []rune, pos int, _ readline.DelayedTabContext) (string, []*readline.CompletionGroup) {
		ctx := string(line)
		var completions []string
		
		compGroup := []*readline.CompletionGroup{
			&readline.CompletionGroup{
				TrimSlash: false,
				NoSpace: true,
			},
		}

		ctx = strings.TrimLeft(ctx, " ")
		if len(ctx) == 0 {
			return "", compGroup
		}

		fields := strings.Split(ctx, " ")
		if len(fields) == 0 {
			return "", compGroup
		}
		query := fields[len(fields) - 1]

		ctx = aliases.Resolve(ctx)
		
		if len(fields) == 1 {
			fileCompletions := fileComplete(query, ctx, fields)
			if len(fileCompletions) != 0 {
				for _, f := range fileCompletions {
					name := strings.Replace(query + f, "~", curuser.HomeDir, 1)
					if info, err := os.Stat(name); err == nil && info.Mode().Perm() & 0100 == 0 {
						continue
					}
					completions = append(completions, f)
				}
				compGroup[0].Suggestions = completions
				return "", compGroup
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
			
			compGroup[0].Suggestions = completions
			return query, compGroup
		} else {
			if completecb, ok := luaCompletions["command." + fields[0]]; ok {
				err := l.CallByParam(lua.P{
					Fn: completecb,
					NRet: 1,
					Protect: true,
				})

				if err != nil {
					return "", compGroup
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
												completions = append(completions, fileComplete(query, ctx, fields)...)
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
				completions = fileComplete(query, ctx, fields)
			}
		}

		compGroup[0].Suggestions = completions
		return "", compGroup
	}

	return &lineReader{
		rl,
	}
}

func (lr *lineReader) Read() (string, error) {
	hooks.Em.Emit("command.precmd", nil)
	s, err := lr.rl.Readline()
	// this is so dumb
	if err == readline.EOF {
		return "", io.EOF
	}

	return s, err // might get another error
}

func (lr *lineReader) SetPrompt(prompt string) {
	halfPrompt := strings.Split(prompt, "\n")
	if len(halfPrompt) > 1 {
		lr.rl.SetPrompt(strings.Join(halfPrompt[:len(halfPrompt) - 1], "\n"))
		lr.rl.MultilinePrompt = halfPrompt[len(halfPrompt) - 1:][0]
	} else {
		// print cursor up ansi code
		fmt.Printf("\033[1A")
		lr.rl.SetPrompt("")
		lr.rl.MultilinePrompt = halfPrompt[len(halfPrompt) - 1:][0]
	}
}

func (lr *lineReader) AddHistory(cmd string) {
	return
}

func (lr *lineReader) ClearInput() {
	return
}

func (lr *lineReader) Resize() {
	return
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

	mod := l.SetFuncs(l.NewTable(), lrLua)

	return mod
}

func (lr *lineReader) luaAddHistory(l *lua.LState) int {
	cmd := l.CheckString(1)
	lr.AddHistory(cmd)

	return 0
}

func (lr *lineReader) luaSize(l *lua.LState) int {
	return 0
}

func (lr *lineReader) luaGetHistory(l *lua.LState) int {
	return 0
}

func (lr *lineReader) luaAllHistory(l *lua.LState) int {
	return 0
}

func (lr *lineReader) luaClearHistory(l *lua.LState) int {
	return 0
}
