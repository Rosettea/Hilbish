package main

import (
	"plugin"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

// #interface module
// native module loading
/* The hilbish.module interface provides a function
to load Hilbish plugins/modules.
Hilbish modules are Go-written plugins (see https://pkg.go.dev/plugin)
that are used to add functionality to Hilbish that cannot be written
n Lua for any reason.

To make a valid native module, the Go plugin
has to export a Loader function with a signature like so:
`func(*rt.Runtime) rt.Value`
`rt` in this case refers to the Runtime type at
https://pkg.go.dev/github.com/arnodel/golua@master/runtime#Runtime
Hilbish uses this package as its Lua runtime. You will need to read
it to use it for a native plugin.
*/
func moduleLoader(rtm *rt.Runtime) *rt.Table {
	exports := map[string]util.LuaExport{
		"load": {moduleLoad, 2, false},
	}

	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	return mod
}

// #interface module
// load(path)
// Loads a module at the designated `path`.
// It will throw if any error occurs.
func moduleLoad(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(1); err != nil {
		return nil, err
	}
	
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	value, err := p.Lookup("Loader")
	if err != nil {
		return nil, err
	}

	loader, ok := value.(func(*rt.Runtime) rt.Value)
	if !ok {
		return nil, nil
	}

	val := loader(t.Runtime)

	return c.PushingNext1(t.Runtime, val), nil
}
