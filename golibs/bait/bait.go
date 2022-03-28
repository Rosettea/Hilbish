package bait

import (
	"fmt"
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib/packagelib"
	"github.com/chuckpreslar/emission"
)

type Bait struct{
	Em *emission.Emitter
	Loader packagelib.Loader
}

func New() Bait {
	emitter := emission.NewEmitter()
	emitter.RecoverWith(func(hookname, hookfunc interface{}, err error) {
		emitter.Off(hookname, hookfunc)
		fmt.Println(err)
	})
	b := Bait{
		Em: emitter,
	}
	b.Loader = packagelib.Loader{
		Load: b.LoaderFunc,
		Name: "bait",
	}

	return b
}

func (b *Bait) LoaderFunc(rtm *rt.Runtime) (rt.Value, func()) {
	exports := map[string]util.LuaExport{
		"catch": util.LuaExport{b.bcatch, 2, false},
		"catchOnce": util.LuaExport{b.bcatchOnce, 2, false},
		"throw": util.LuaExport{b.bthrow, 1, true},
	}
	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

/*
	util.Document(L, mod,
`Bait is the event emitter for Hilbish. Why name it bait?
Because it throws hooks that you can catch (emits events
that you can listen to) and because why not, fun naming
is fun. This is what you will use if you want to listen
in on hooks to know when certain things have happened,
like when you've changed directory, a command has
failed, etc. To find all available hooks, see doc hooks.`)
*/

	return rt.TableValue(mod), nil
}

func handleHook(t *rt.Thread, c *rt.GoCont, name string, catcher *rt.Closure, args ...interface{}) {
	funcVal := rt.FunctionValue(catcher)
	var luaArgs []rt.Value
	for _, arg := range args {
		var luarg rt.Value
		switch arg.(type) {
			case rt.Value: luarg = arg.(rt.Value)
			default: luarg = rt.AsValue(arg)
		}
		luaArgs = append(luaArgs, luarg)
	}
	_, err := rt.Call1(t, funcVal, luaArgs...)
	if err != nil {
		e := rt.NewError(rt.StringValue(err.Error()))
		e = e.AddContext(c.Next(), 1)
		// panicking here won't actually cause hilbish to panic and instead will
		// print the error and remove the hook (look at emission recover from above)
		panic(e)
	}
}

// throw(name, ...args)
// Throws a hook with `name` with the provided `args`
// --- @param name string
// --- @vararg any
func (b *Bait) bthrow(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	name, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	ifaceSlice := make([]interface{}, len(c.Etc()))
	for i, v := range c.Etc() {
		ifaceSlice[i] = v
	}
	b.Em.Emit(name, ifaceSlice...)

	return c.Next(), nil
}

// catch(name, cb)
// Catches a hook with `name`. Runs the `cb` when it is thrown
// --- @param name string
// --- @param cb function
func (b *Bait) bcatch(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	name, catcher, err := util.HandleStrCallback(t, c)
	if err != nil {
		return nil, err
	}

	b.Em.On(name, func(args ...interface{}) {
		handleHook(t, c, name, catcher, args...)
	})

	return c.Next(), nil
}

// catchOnce(name, cb)
// Same as catch, but only runs the `cb` once and then removes the hook
// --- @param name string
// --- @param cb function
func (b *Bait) bcatchOnce(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	name, catcher, err := util.HandleStrCallback(t, c)
	if err != nil {
		return nil, err
	}

	b.Em.Once(name, func(args ...interface{}) {
		handleHook(t, c, name, catcher, args...)
	})

	return c.Next(), nil
}
