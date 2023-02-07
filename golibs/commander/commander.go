// library for custom commands
/*
Commander is a library for writing custom commands in Lua.
In order to make it easier to write commands for Hilbish,
not require separate scripts and to be able to use in a config,
the Commander library exists. This is like a very simple wrapper
that works with Hilbish for writing commands. Example:

```lua
local commander = require 'commander'

commander.register('hello', function(args, sinks)
	sinks.out:writeln 'Hello world!'
end)
```

In this example, a command with the name of `hello` is created
that will print `Hello world!` to output. One question you may
have is: What is the `sinks` parameter?

The `sinks` parameter is a table with 3 keys: `in`, `out`,
and `err`. The values of these is a @Sink.

- `in` is the standard input. You can read from this sink
to get user input. (**This is currently unimplemented.**)
- `out` is standard output. This is usually where text meant for
output should go.
- `err` is standard error. This sink is for writing errors, as the
name would suggest.
*/
package commander

import (
	"hilbish/util"
	"hilbish/golibs/bait"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib/packagelib"
)

type Commander struct{
	Events *bait.Bait
	Loader packagelib.Loader
}

func New(rtm *rt.Runtime) Commander {
	c := Commander{
		Events: bait.New(rtm),
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
