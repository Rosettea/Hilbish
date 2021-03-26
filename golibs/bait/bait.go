package bait

import (
	"github.com/chuckpreslar/emission"
	"github.com/yuin/gopher-lua"
)

type Bait struct{}

func New() Bait {
	return Bait{}
}

func (b *Bait) Loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{}
	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)

	return 1
}

