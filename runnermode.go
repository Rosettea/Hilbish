package main

import (
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

func runnerModeLoader(rtm *rt.Runtime) *rt.Table {
	exports := map[string]util.LuaExport{
		"sh": {shRunner, 1, false},
		"lua": {luaRunner, 1, false},
		"setMode": {hlrunnerMode, 1, false},
	}

	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	return mod
}

func shRunner(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	input, exitCode, err := handleSh(cmd)
	var luaErr rt.Value = rt.NilValue
	if err != nil {
		luaErr = rt.StringValue(err.Error())
	}

	return c.PushingNext(t.Runtime, rt.StringValue(input), rt.IntValue(int64(exitCode)), luaErr), nil
}

func luaRunner(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	input, exitCode, err := handleLua(cmd)
	var luaErr rt.Value = rt.NilValue
	if err != nil {
		luaErr = rt.StringValue(err.Error())
	}

	return c.PushingNext(t.Runtime, rt.StringValue(input), rt.IntValue(int64(exitCode)), luaErr), nil
}
