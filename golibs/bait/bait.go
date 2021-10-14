package bait

import (
	"github.com/chuckpreslar/emission"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

type Bait struct{
	Em *emission.Emitter
}

func New() Bait {
	return Bait{
		Em: emission.NewEmitter(),
	}
}

func (b *Bait) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{})
	L.SetField(mod, "throw", luar.New(L, b.throw))
	L.SetField(mod, "catch", luar.New(L, b.catch))

	L.Push(mod)

	return 1
}

func (b *Bait) throw(name string, args ...interface{}) {
	b.Em.Emit(name, args...)
}

func (b *Bait) catch(name string, catcher func(...interface{})) {
	b.Em.On(name, catcher)
}
