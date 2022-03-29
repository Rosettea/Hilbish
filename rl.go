package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/maxlandon/readline"
	rt "github.com/arnodel/golua/runtime"
)

type lineReader struct {
	rl *readline.Instance
}
var fileHist *fileHistory
/*
var hinter lua.LValue = lua.LNil
var highlighter lua.LValue = lua.LNil
*/

// other gophers might hate this naming but this is local, shut up
func newLineReader(prompt string, noHist bool) *lineReader {
	rl := readline.NewInstance()
	// we don't mind hilbish.read rl instances having completion,
	// but it cant have shared history
	if !noHist {
		fileHist = newFileHistory()
		rl.SetHistoryCtrlR("History", fileHist)
		rl.HistoryAutoWrite = false
	}
	rl.ShowVimMode = false
	rl.ViModeCallback = func(mode readline.ViMode) {
		modeStr := ""
		switch mode {
			case readline.VimKeys: modeStr = "normal"
			case readline.VimInsert: modeStr = "insert"
			case readline.VimDelete: modeStr = "delete"
			case readline.VimReplaceOnce:
			case readline.VimReplaceMany: modeStr = "replace"
		}
		setVimMode(modeStr)
	}
	rl.ViActionCallback = func(action readline.ViAction, args []string) {
		actionStr := ""
		switch action {
			case readline.VimActionPaste: actionStr = "paste"
			case readline.VimActionYank: actionStr = "yank"
		}
		hooks.Em.Emit("hilbish.vimAction", actionStr, args)
	}
	/*
	rl.HintText = func(line []rune, pos int) []rune {
		if hinter == lua.LNil {
			return []rune{}
		}

		err := l.CallByParam(lua.P{
			Fn: hinter,
			NRet: 1,
			Protect: true,
		}, lua.LString(string(line)), lua.LNumber(pos))
		if err != nil {
			fmt.Println(err)
			return []rune{}
		}
		
		retVal := l.Get(-1)
		hintText := ""
		if luaStr, ok := retVal.(lua.LString); retVal != lua.LNil && ok {
			hintText = luaStr.String()
		}
		
		return []rune(hintText)
	}
	rl.SyntaxHighlighter = func(line []rune) string {
		if highlighter == lua.LNil {
			return string(line)
		}
		err := l.CallByParam(lua.P{
			Fn: highlighter,
			NRet: 1,
			Protect: true,
		}, lua.LString(string(line)))
		if err != nil {
			fmt.Println(err)
			return string(line)
		}
		
		retVal := l.Get(-1)
		highlighted := ""
		if luaStr, ok := retVal.(lua.LString); retVal != lua.LNil && ok {
			highlighted = luaStr.String()
		}
		
		return highlighted
	}
	*/
	rl.TabCompleter = func(line []rune, pos int, _ readline.DelayedTabContext) (string, []*readline.CompletionGroup) {
		ctx := string(line)
		var completions []string

		var compGroup []*readline.CompletionGroup

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
			completions, prefix := binaryComplete(query, ctx, fields)

			compGroup = append(compGroup, &readline.CompletionGroup{
				TrimSlash: false,
				NoSpace: true,
				Suggestions: completions,
			})

			return prefix, compGroup
		} else {
			if completecb, ok := luaCompletions["command." + fields[0]]; ok {
				luaFields := rt.NewTable()
				for i, f := range fields {
					luaFields.Set(rt.IntValue(int64(i + 1)), rt.StringValue(f))
				}

				// we must keep the holy 80 cols
				luacompleteTable, err := rt.Call1(l.MainThread(), 
				rt.FunctionValue(completecb), rt.StringValue(query),
				rt.StringValue(ctx), rt.TableValue(luaFields))

				if err != nil {
					return "", compGroup
				}

				/*
					as an example with git,
					completion table should be structured like:
					{
						{
							items = {
								'add',
								'clone',
								'init'
							},
							type = 'grid'
						},
						{
							items = {
								'-c',
								'--git-dir'
							},
							type = 'list'
						}
					}
					^ a table of completion groups.
					it is the responsibility of the completer
					to work on subcommands and subcompletions
				*/
				if cmpTbl, ok := luacompleteTable.TryTable(); ok {
					nextVal := rt.NilValue
					for {
						next, val, ok := cmpTbl.Next(nextVal)
						if next == rt.NilValue {
							break
						}
						nextVal = next

						_, ok = next.TryInt()
						valTbl, okk := val.TryTable()
						if !ok || !okk {
							// TODO: error?
							break
						}

						luaCompType := valTbl.Get(rt.StringValue("type"))
						luaCompItems := valTbl.Get(rt.StringValue("items"))

						compType, ok := luaCompType.TryString()
						compItems, okk := luaCompItems.TryTable()
						if !ok || !okk {
							// TODO: error
							break
						}

						var items []string
						itemDescriptions := make(map[string]string)
						nxVal := rt.NilValue
						for {
							nx, vl, _ := compItems.Next(nxVal)
							if nx == rt.NilValue {
								break
							}
							nxVal = nx

							if tstr := nx.Type(); tstr == rt.StringType {
								// ['--flag'] = {'description', '--flag-alias'}
								nxStr, ok := nx.TryString()
								vlTbl, okk := vl.TryTable()
								if !ok || !okk {
									// TODO: error
									continue
								}
								items = append(items, nxStr)
								itemDescription, ok := vlTbl.Get(rt.IntValue(1)).TryString()
								if !ok {
									// TODO: error
									continue
								}
								itemDescriptions[nxStr] = itemDescription
							} else if tstr == rt.IntType {
								vlStr, okk := vl.TryString()
								if !okk {
									// TODO: error
									continue
								}
								items = append(items, vlStr)
							} else {
								// TODO: error
								continue
							}
						}

						var dispType readline.TabDisplayType
						switch compType {
							case "grid": dispType = readline.TabDisplayGrid
							case "list": dispType = readline.TabDisplayList
							// need special cases, will implement later
							//case "map": dispType = readline.TabDisplayMap
						}

						compGroup = append(compGroup, &readline.CompletionGroup{
							DisplayType: dispType,
							Descriptions: itemDescriptions,
							Suggestions: items,
							TrimSlash: false,
							NoSpace: true,
						})
					}
				}
			}

			if len(compGroup) == 0 {
				completions = fileComplete(query, ctx, fields)
				compGroup = append(compGroup, &readline.CompletionGroup{
					TrimSlash: false,
					NoSpace: true,
					Suggestions: completions,
				})
			}
		}
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
		fmt.Println("")
		return "", io.EOF
	}

	return s, err // might get another error
}

func (lr *lineReader) SetPrompt(p string) {
	halfPrompt := strings.Split(p, "\n")
	if len(halfPrompt) > 1 {
		lr.rl.Multiline = true
		lr.rl.SetPrompt(strings.Join(halfPrompt[:len(halfPrompt) - 1], "\n"))
		lr.rl.MultilinePrompt = halfPrompt[len(halfPrompt) - 1:][0]
	} else {
		lr.rl.Multiline = false
		lr.rl.MultilinePrompt = ""
		lr.rl.SetPrompt(p)
	}
	if initialized && !running {
		lr.rl.RefreshPromptInPlace("")
	}
}

func (lr *lineReader) AddHistory(cmd string) {
	fileHist.Write(cmd)
}

func (lr *lineReader) ClearInput() {
	return
}

func (lr *lineReader) Resize() {
	return
}

// lua module
/*
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

func (lr *lineReader) luaSize(L *lua.LState) int {
	L.Push(lua.LNumber(fileHist.Len()))

	return 1
}

func (lr *lineReader) luaGetHistory(L *lua.LState) int {
	idx := L.CheckInt(1)
	cmd, _ := fileHist.GetLine(idx)
	L.Push(lua.LString(cmd))

	return 1
}

func (lr *lineReader) luaAllHistory(L *lua.LState) int {
	tbl := L.NewTable()
	size := fileHist.Len()

	for i := 1; i < size; i++ {
		cmd, _ := fileHist.GetLine(i)
		tbl.Append(lua.LString(cmd))
	}

	L.Push(tbl)

	return 1
}

func (lr *lineReader) luaClearHistory(l *lua.LState) int {
	fileHist.clear()
	return 0
}
*/
