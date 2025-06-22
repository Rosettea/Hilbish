// the event emitter
/*
Bait is the event emitter for Hilbish. Much like Node.js and
its `events` system, many actions in Hilbish emit events.
Unlike Node.js, Hilbish events are global. So make sure to
pick a unique name!<nl>
Usage of the Bait module consists of userstanding
event-driven architecture, but it's pretty simple:
If you want to act on a certain event, you can `catch` it.
You can act on events via callback functions.<nl>
Examples of this are in the Hilbish default config!
Consider this part of it:<nl>

```lua
bait.catch('command.exit', function(code)
	running = false
	doPrompt(code ~= 0)
	doNotifyPrompt()
end)
```

What this does is, whenever the `command.exit` event is thrown,
this function will set the user prompt.
*/
package bait

import (
	"errors"

	"hilbish/util"

	"github.com/arnodel/golua/lib/packagelib"
	rt "github.com/arnodel/golua/runtime"
)

type listenerType int

const (
	goListener listenerType = iota
	luaListener
)

// Recoverer is a function which is called when a panic occurs in an event.
type Recoverer func(event string, handler *Listener, err interface{})

// Listener is a struct that holds the handler for an event.
type Listener struct {
	typ       listenerType
	once      bool
	caller    func(...interface{}) rt.Value
	luaCaller *rt.Closure
}

type Bait struct {
	Loader    packagelib.Loader
	recoverer Recoverer
	handlers  map[string][]*Listener
	rtm       *rt.Runtime
}

// New creates a new Bait instance.
func New(rtm *rt.Runtime) *Bait {
	b := &Bait{
		handlers: make(map[string][]*Listener),
		rtm:      rtm,
	}
	b.Loader = packagelib.Loader{
		Load: b.loaderFunc,
		Name: "bait",
	}

	return b
}

// Emit throws an event.
func (b *Bait) Emit(event string, args ...interface{}) []rt.Value {
	var returns []rt.Value
	handles := b.handlers[event]
	if handles == nil {
		return nil
	}

	for idx, handle := range handles {
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
				case rt.Value:
					luarg = arg.(rt.Value)
				default:
					luarg = rt.AsValue(arg)
				}
				luaArgs = append(luaArgs, luarg)
			}
			luaRet, err := rt.Call1(b.rtm.MainThread(), funcVal, luaArgs...)
			if err != nil {
				if event != "error" {
					b.Emit("error", event, handle.luaCaller, err.Error())
					return nil
				}
				// if there is an error in an error event handler, panic instead
				// (calls the go recoverer function)
				panic(err)
			}

			if luaRet != rt.NilValue {
				returns = append(returns, luaRet)
			}
		} else {
			ret := handle.caller(args...)
			if ret != rt.NilValue {
				returns = append(returns, ret)
			}
		}

		if handle.once {
			b.removeListener(event, idx)
		}
	}

	return returns
}

// On adds a Go function handler for an event.
func (b *Bait) On(event string, handler func(...interface{}) rt.Value) *Listener {
	listener := &Listener{
		typ:    goListener,
		caller: handler,
	}

	b.addListener(event, listener)
	return listener
}

// OnLua adds a Lua function handler for an event.
func (b *Bait) OnLua(event string, handler *rt.Closure) *Listener {
	listener := &Listener{
		typ:       luaListener,
		luaCaller: handler,
	}
	b.addListener(event, listener)

	return listener
}

// Off removes a Go function handler for an event.
func (b *Bait) Off(event string, listener *Listener) {
	handles := b.handlers[event]

	for i, handle := range handles {
		if handle == listener {
			b.removeListener(event, i)
		}
	}
}

// OffLua removes a Lua function handler for an event.
func (b *Bait) OffLua(event string, handler *rt.Closure) {
	handles := b.handlers[event]

	for i, handle := range handles {
		if handle.luaCaller == handler {
			b.removeListener(event, i)
		}
	}
}

// Once adds a Go function listener for an event that only runs once.
func (b *Bait) Once(event string, handler func(...interface{}) rt.Value) *Listener {
	listener := &Listener{
		typ:    goListener,
		once:   true,
		caller: handler,
	}
	b.addListener(event, listener)

	return listener
}

// OnceLua adds a Lua function listener for an event that only runs once.
func (b *Bait) OnceLua(event string, handler *rt.Closure) *Listener {
	listener := &Listener{
		typ:       luaListener,
		once:      true,
		luaCaller: handler,
	}
	b.addListener(event, listener)

	return listener
}

// SetRecoverer sets the function to be executed when a panic occurs in an event.
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
	b.handlers[event][idx] = b.handlers[event][len(b.handlers[event])-1]

	b.handlers[event] = b.handlers[event][:len(b.handlers[event])-1]
}

func (b *Bait) callRecoverer(event string, handler *Listener, err interface{}) {
	if b.recoverer == nil {
		panic(err)
	}
	b.recoverer(event, handler, err)
}

func (b *Bait) loaderFunc(rtm *rt.Runtime) (rt.Value, func()) {
	exports := map[string]util.LuaExport{
		"catch":     util.LuaExport{b.bcatch, 2, false},
		"catchOnce": util.LuaExport{b.bcatchOnce, 2, false},
		"throw":     util.LuaExport{b.bthrow, 1, true},
		"release":   util.LuaExport{b.brelease, 2, false},
		"hooks":     util.LuaExport{b.bhooks, 1, false},
	}
	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	return rt.TableValue(mod), nil
}

// catch(name, cb)
// Catches an event. This function can be used to act on events.
// #param name string The name of the hook.
// #param cb function The function that will be called when the hook is thrown.
/*
#example
bait.catch('hilbish.exit', function()
	print 'Goodbye Hilbish!'
end)
#example
*/
func (b *Bait) bcatch(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	name, catcher, err := util.HandleStrCallback(t, c)
	if err != nil {
		return nil, err
	}

	b.OnLua(name, catcher)

	return c.Next(), nil
}

// catchOnce(name, cb)
// Catches an event, but only once. This will remove the hook immediately after it runs for the first time.
// #param name string The name of the event
// #param cb function The function that will be called when the event is thrown.
func (b *Bait) bcatchOnce(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	name, catcher, err := util.HandleStrCallback(t, c)
	if err != nil {
		return nil, err
	}

	b.OnceLua(name, catcher)

	return c.Next(), nil
}

// hooks(name) -> table
// Returns a table of functions that are hooked on an event with the corresponding `name`.
// #param name string The name of the hook
// #returns table<function>
func (b *Bait) bhooks(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	evName, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	noHooks := errors.New("no hooks for event " + evName)

	handlers := b.handlers[evName]
	if handlers == nil {
		return nil, noHooks
	}

	luaHandlers := rt.NewTable()
	for _, handler := range handlers {
		if handler.typ != luaListener {
			continue
		}
		luaHandlers.Set(rt.IntValue(luaHandlers.Len()+1), rt.FunctionValue(handler.luaCaller))
	}

	if luaHandlers.Len() == 0 {
		return nil, noHooks
	}

	return c.PushingNext1(t.Runtime, rt.TableValue(luaHandlers)), nil
}

// release(name, catcher)
// Removes the `catcher` for the event with `name`.
// For this to work, `catcher` has to be the same function used to catch
// an event, like one saved to a variable.
// #param name string Name of the event the hook is on
// #param catcher function Hook function to remove
/*
#example
local hookCallback = function() print 'hi' end

bait.catch('event', hookCallback)

-- a little while later....
bait.release('event', hookCallback)
-- and now hookCallback will no longer be ran for the event.
#example
*/
func (b *Bait) brelease(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	name, catcher, err := util.HandleStrCallback(t, c)
	if err != nil {
		return nil, err
	}

	b.OffLua(name, catcher)

	return c.Next(), nil
}

// throw(name, ...args)
// #param name string The name of the hook.
// #param args ...any The arguments to pass to the hook.
// Throws a hook with `name` with the provided `args`.
/*
#example
bait.throw('greeting', 'world')

-- This can then be listened to via
bait.catch('gretting', function(greetTo)
	print('Hello ' .. greetTo)
end)
#example
*/
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
	ret := b.Emit(name, ifaceSlice...)

	return c.PushingNext(t.Runtime, ret...), nil
}
