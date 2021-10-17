package terminal

import (
	"os"

	"hilbish/util"

	"golang.org/x/term"
	"github.com/yuin/gopher-lua"
)

var termState *term.State

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	util.Document(L, mod, "The terminal library is a simple and lower level library for certain terminal interactions.")

	L.Push(mod)

	return 1
}

var exports = map[string]lua.LGFunction{
	"setRaw": termraw,
	"restoreState": termrestoreState,
	"size": termsize,
	"saveState": termsaveState,
}

// size()
// Gets the dimensions of the terminal. Returns a table with `width` and `height`
// Note: this is not the size in relation to the dimensions of the display
func termsize(L *lua.LState) int {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	dimensions := L.NewTable()
	L.SetField(dimensions, "width", lua.LNumber(w))
	L.SetField(dimensions, "height", lua.LNumber(h))

	L.Push(dimensions)
	return 1
}

// saveState()
// Saves the current state of the terminal
func termsaveState(L *lua.LState) int {
	state, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	termState = state
	return 0
}

// restoreState()
// Restores the last saved state of the terminal
func termrestoreState(L *lua.LState) int {
	err := term.Restore(int(os.Stdin.Fd()), termState)
	if err != nil {
		L.RaiseError(err.Error())
	}

	return 0
}

// setRaw()
// Puts the terminal in raw mode
func termraw(L *lua.LState) int {
	_, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		L.RaiseError(err.Error())
	}

	return 0
}
