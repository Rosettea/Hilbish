// +build goreadline,!hilbiline

package main

// Here we define a generic interface for readline and hilbiline,
// making them interchangable during build time
// this is hilbiline's, as is obvious by the filename

import (
	"io"
	"strings"
	"fmt"

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

	return &lineReader{
		rl,
	}
}

func (lr *lineReader) Read() (string, error) {
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
func (lr *lineReader) Loader() *lua.LTable {
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
