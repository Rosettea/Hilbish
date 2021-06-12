package commander

import (
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
	var exports = map[string]lua.LGFunction{
		"register": c.register,
		"deregister": c.deregister,
	}
	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)

	return 1
}

func (c *Commander) register(L *lua.LState) int {
	cmdName := L.CheckString(1)
	cmd := L.CheckFunction(2)

	c.Events.Emit("commandRegister", cmdName, cmd)

	return 0
}

func (c *Commander) deregister(L *lua.LState) int {
	cmdName := L.CheckString(1)

	c.Events.Emit("commandDeregister", cmdName)

	return 0
}
