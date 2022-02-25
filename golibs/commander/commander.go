package commander

import (
	"hilbish/util"

	"github.com/chuckpreslar/emission"
	"github.com/yuin/gopher-lua"
)

type Commander struct{
	Events *emission.Emitter
}

func New() Commander {
	return Commander{
		Events: emission.NewEmitter(),
	}
}

func (c *Commander) Loader(L *lua.LState) int {
	exports := map[string]lua.LGFunction{
		"register": c.cregister,
		"deregister": c.cderegister,
	}
	mod := L.SetFuncs(L.NewTable(), exports)
	util.Document(L, mod, "Commander is Hilbish's custom command library, a way to write commands in Lua.")
	L.Push(mod)

	return 1
}

// register(name, cb)
// Register a command with `name` that runs `cb` when ran
func (c *Commander) cregister(L *lua.LState) int {
	cmdName := L.CheckString(1)
	cmd := L.CheckFunction(2)

	c.Events.Emit("commandRegister", cmdName, cmd)

	return 0
}

// deregister(name)
// Deregisters any command registered with `name`
func (c *Commander) cderegister(L *lua.LState) int {
	cmdName := L.CheckString(1)

	c.Events.Emit("commandDeregister", cmdName)

	return 0
}
