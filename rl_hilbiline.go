// +build hilbiline,!goreadline

package main

// Here we define a generic interface for readline and hilbiline,
// making them interchangable during build time
// this is hilbiline's, as is obvious by the filename

import (
	"github.com/Rosettea/Hilbiline"
	"github.com/yuin/gopher-lua"
)

type lineReader struct {
	hl *hilbiline.HilbilineState
}

// other gophers might hate this naming but this is local, shut up
func newLineReader(prompt string) *lineReader {
	hl := hilbiline.New(prompt)

	return &lineReader{
		&hl,
	}
}

func (lr *lineReader) Read() (string, error) {
	return lr.hl.Read()
}

func (lr *lineReader) SetPrompt(prompt string) {
	lr.hl.SetPrompt(prompt)
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
