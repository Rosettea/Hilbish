package main

import (
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

// #interface runner
// interactive command runner customization
// The runner interface contains functions that allow the user to change
// how Hilbish interprets interactive input.
// Users can add and change the default runner for interactive input to any
// language or script of their choosing. A good example is using it to
// write command in Fennel.
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

// #interface runner
// setMode(cb)
// This is the same as the `hilbish.runnerMode` function. It takes a callback,
// which will be used to execute all interactive input.
// In normal cases, neither callbacks should be overrided by the user,
// as the higher level functions listed below this will handle it.
func _runnerMode() {}

// #interface runner
// sh(cmd)
// Runs a command in Hilbish's shell script interpreter.
// This is the equivalent of using `source`.
func shRunner(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	_, exitCode, cont, err := execSh(aliases.Resolve(cmd))
	var luaErr rt.Value = rt.NilValue
	if err != nil {
		luaErr = rt.StringValue(err.Error())
	}
	runnerRet := rt.NewTable()
	runnerRet.Set(rt.StringValue("input"), rt.StringValue(cmd))
	runnerRet.Set(rt.StringValue("exitCode"), rt.IntValue(int64(exitCode)))
	runnerRet.Set(rt.StringValue("continue"), rt.BoolValue(cont))
	runnerRet.Set(rt.StringValue("err"), luaErr)

	return c.PushingNext(t.Runtime, rt.TableValue(runnerRet)), nil
}

// #interface runner
// lua(cmd)
// Evaluates `cmd` as Lua input. This is the same as using `dofile`
// or `load`, but is appropriated for the runner interface.
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
	runnerRet := rt.NewTable()
	runnerRet.Set(rt.StringValue("input"), rt.StringValue(input))
	runnerRet.Set(rt.StringValue("exitCode"), rt.IntValue(int64(exitCode)))
	runnerRet.Set(rt.StringValue("err"), luaErr)

	return c.PushingNext(t.Runtime, rt.TableValue(runnerRet)), nil
}
