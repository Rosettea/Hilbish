package commander

import (
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib/packagelib"
	"github.com/chuckpreslar/emission"
)

type Commander struct{
	Events *emission.Emitter
	Loader packagelib.Loader
}

func New() Commander {
	c := Commander{
		Events: emission.NewEmitter(),
	}
	c.Loader = packagelib.Loader{
		Load: c.loaderFunc,
		Name: "commander",
	}

	return c
}

func (c *Commander) loaderFunc(rtm *rt.Runtime) (rt.Value, func()) {
	exports := map[string]util.LuaExport{
		"register": util.LuaExport{c.cregister, 2, false},
		"deregister": util.LuaExport{c.cderegister, 1, false},
	}
	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)
	util.Document(mod, "Commander is Hilbish's custom command library, a way to write commands in Lua.")

	return rt.TableValue(mod), nil
}

// register(name, cb)
// Register a command with `name` that runs `cb` when ran
// --- @param name string
// --- @param cb function
func (c *Commander) cregister(t *rt.Thread, ct *rt.GoCont) (rt.Cont, error) {
	cmdName, cmd, err := util.HandleStrCallback(t, ct)
	if err != nil {
		return nil, err
	}

	c.Events.Emit("commandRegister", cmdName, cmd)

	return ct.Next(), err
}

// deregister(name)
// Deregisters any command registered with `name`
// --- @param name string
func (c *Commander) cderegister(t *rt.Thread, ct *rt.GoCont) (rt.Cont, error) {
	if err := ct.Check1Arg(); err != nil {
		return nil, err
	}
	cmdName, err := ct.StringArg(0)
	if err != nil {
		return nil, err
	}

	c.Events.Emit("commandDeregister", cmdName)

	return ct.Next(), err
}
