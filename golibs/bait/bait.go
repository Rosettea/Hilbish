package bait

import (
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib/packagelib"
)

type listenerType int
const (
	goListener listenerType = iota
	luaListener
)

type Recoverer func(event string, handler *Listener, err interface{})

type Listener struct{
	typ listenerType
	once bool
	caller func(...interface{})
	luaCaller *rt.Closure
}

type Bait struct{
	Loader packagelib.Loader
	recoverer Recoverer
	handlers map[string][]*Listener
	rtm *rt.Runtime
}

func New(rtm *rt.Runtime) Bait {
	b := Bait{
		handlers: make(map[string][]*Listener),
		rtm: rtm,
	}
	b.Loader = packagelib.Loader{
		Load: b.loaderFunc,
		Name: "bait",
	}

	return b
}

func (b *Bait) Emit(event string, args ...interface{}) {
	handles := b.handlers[event]
	if handles == nil {
		return
	}

	for _, handle := range handles {
		defer func() {
			if err := recover(); err != nil {
				b.callRecoverer(event, handle, err)
			}
		}()

		if handle.typ == luaListener {
			funcVal := rt.FunctionValue(handle.luaCaller)
			var luaArgs []rt.Value
			for _, arg := range args {
				var luarg rt.Value
				switch arg.(type) {
					case rt.Value: luarg = arg.(rt.Value)
					default: luarg = rt.AsValue(arg)
				}
				luaArgs = append(luaArgs, luarg)
			}
			_, err := rt.Call1(b.rtm.MainThread(), funcVal, luaArgs...)
			if err != nil {
				// panicking here won't actually cause hilbish to panic and instead will
				// print the error and remove the hook. reference the recoverer function in lua.go
				panic(err)
			}
		} else {
			handle.caller(args...)
		}

		if handle.once {
			b.Off(event, handle)
		}
	}
}

func (b *Bait) On(event string, handler func(...interface{})) *Listener {
	listener := &Listener{
		typ: goListener,
		caller: handler,
	}

	b.addListener(event, listener)
	return listener
}

func (b *Bait) OnLua(event string, handler *rt.Closure) *Listener {
	listener :=&Listener{
		typ: luaListener,
		luaCaller: handler,
	}
	b.addListener(event, listener)

	return listener
}

func (b *Bait) Off(event string, listener *Listener) {
	handles := b.handlers[event]

	for i, handle := range handles {
		if handle == listener {
			b.removeListener(event, i)
		}
	}
}

func (b *Bait) OffLua(event string, handler *rt.Closure) {
	handles := b.handlers[event]

	for i, handle := range handles {
		if handle.luaCaller == handler {
			b.removeListener(event, i)
		}
	}
}

func (b *Bait) Once(event string, handler func(...interface{})) *Listener {
	listener := &Listener{
		typ: goListener,
		once: true,
		caller: handler,
	}
	b.addListener(event, listener)

	return listener
}

func (b *Bait) OnceLua(event string, handler *rt.Closure) *Listener {
	listener := &Listener{
		typ: luaListener,
		once: true,
		luaCaller: handler,
	}
	b.addListener(event, listener)

	return listener
}

func (b *Bait) SetRecoverer(recoverer Recoverer) {
	b.recoverer = recoverer
}

func (b *Bait) addListener(event string, listener *Listener) {
	if b.handlers[event] == nil {
		b.handlers[event] = []*Listener{}
	}

	b.handlers[event] = append(b.handlers[event], listener)
}


func (b *Bait) removeListener(event string, idx int) {
	b.handlers[event][idx] = b.handlers[event][len(b.handlers[event]) - 1]

	b.handlers[event] = b.handlers[event][:len(b.handlers[event]) - 1]
}

func (b *Bait) callRecoverer(event string, handler *Listener, err interface{}) {
	if b.recoverer == nil {
		panic(err)
	}
	b.recoverer(event, handler, err)
}

func (b *Bait) loaderFunc(rtm *rt.Runtime) (rt.Value, func()) {
	exports := map[string]util.LuaExport{
		"catch": util.LuaExport{b.bcatch, 2, false},
		"catchOnce": util.LuaExport{b.bcatchOnce, 2, false},
		"throw": util.LuaExport{b.bthrow, 1, true},
		"release": util.LuaExport{b.brelease, 2, false},
	}
	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	util.Document(mod,
`Bait is the event emitter for Hilbish. Why name it bait?
Because it throws hooks that you can catch (emits events
that you can listen to) and because why not, fun naming
is fun. This is what you will use if you want to listen
in on hooks to know when certain things have happened,
like when you've changed directory, a command has
failed, etc. To find all available hooks, see doc hooks.`)

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
	b.Emit(name, ifaceSlice...)

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

	b.OnLua(name, catcher)

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

	b.OnceLua(name, catcher)

	return c.Next(), nil
}

func (b *Bait) brelease(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	name, catcher, err := util.HandleStrCallback(t, c)
	if err != nil {
		return nil, err
	}

	b.OffLua(name, catcher)

	return c.Next(), nil
}
