package bait

import (
	"hilbish/util"

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
	L.SetField(mod, "throw", luar.New(L, b.bthrow))
	L.SetField(mod, "catch", luar.New(L, b.bcatch))

	util.Document(L, mod, `Bait is the event emitter for Hilbish. Why name it bait?
Because it throws hooks that you can catch
(emits events that you can listen to) and because why not,
fun naming is fun. This is what you will use if you want to
listen in on hooks to know when certain things have happened,
like when you've changed directory, a command has failed, etc.`)
	L.Push(mod)

	return 1
}

// throw(name, ...args)
// Throws a hook with `name` with the provided `args`
func (b *Bait) bthrow(name string, args ...interface{}) {
	b.Em.Emit(name, args...)
}

// catch(name, cb)
// Catches a hook with `name`. Runs the `cb` when it is thrown
func (b *Bait) bcatch(name string, catcher func(...interface{})) {
	b.Em.On(name, catcher)
}
