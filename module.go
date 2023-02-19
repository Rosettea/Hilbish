package main

import (
	"plugin"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

func moduleLoader(rtm *rt.Runtime) *rt.Table {
	exports := map[string]util.LuaExport{
		"load": {moduleLoad, 2, false},
	}

	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	return mod
}

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
