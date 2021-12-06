//go:build !hilbiline
// +build !hilbiline

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

type LineReader struct {
	Prompt string
}

func NewLineReader(prompt string) *LineReader {
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

		for aliases[fields[0]] != "" {
			alias := aliases[fields[0]]
			ctx = alias + strings.TrimPrefix(ctx, fields[0])
			fields = strings.Split(ctx, " ")

			if aliases[fields[0]] == alias {
				break
			}
			if aliases[fields[0]] != "" {
				continue
			}
		}

		if len(fields) == 1 {
			prefixes := []string{"./", "../", "/", "~/"}
			for _, prefix := range prefixes {
				if strings.HasPrefix(fields[0], prefix) {
					fileCompletions := append(completions, readline.FilenameCompleter(query, ctx)...)
					// filter out executables
					for _, f := range fileCompletions {
						name := strings.Replace(f, "~", curuser.HomeDir, 1)
						if info, err := os.Stat(name); err == nil && info.Mode().Perm() & 0100 == 0 {
							continue
						}
						completions = append(completions, f)
					}
					return completions
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

	return &LineReader{
		Prompt: prompt,
	}
}

func (lr *LineReader) Read() (string, error) {
	hooks.Em.Emit("command.precmd", nil)
	return readline.String(lr.Prompt)
}

func (lr *LineReader) SetPrompt(prompt string) {
	lr.Prompt = prompt
}

func (lr *LineReader) AddHistory(cmd string) {
	readline.AddHistory(cmd)
	readline.SaveHistory(defaultHistPath)
}

func (lr *LineReader) ClearInput() {
	readline.ReplaceLine("", 0)
	readline.RefreshLine()
}

func (lr *LineReader) Resize() {
	readline.Resize()
}
