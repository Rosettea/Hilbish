package bait

import (
	"fmt"
	"hilbish/util"

	"github.com/chuckpreslar/emission"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

type Bait struct{
	Em *emission.Emitter
}

func New() Bait {
	emitter := emission.NewEmitter()
	emitter.RecoverWith(func(hookname, hookfunc interface{}, err error) {
		emitter.Off(hookname, hookfunc)
		fmt.Println(err)
	})
	return Bait{
		Em: emitter,
	}
}

func (b *Bait) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{})

	util.Document(L, mod,
`Bait is the event emitter for Hilbish. Why name it bait?
Because it throws hooks that you can catch (emits events
that you can listen to) and because why not, fun naming
is fun. This is what you will use if you want to listen
in on hooks to know when certain things have happened,
like when you've changed directory, a command has
failed, etc. To find all available hooks, see doc hooks.`)

	L.SetField(mod, "throw", luar.New(L, b.bthrow))
	L.SetField(mod, "catch", luar.New(L, b.bcatch))
	L.SetField(mod, "catchOnce", luar.New(L, b.bcatchOnce))

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

// catchOnce(name, cb)
// Same as catch, but only runs the `cb` once and then removes the hook
func (b *Bait) bcatchOnce(name string, catcher func(...interface{})) {
	b.Em.Once(name, catcher)
}
