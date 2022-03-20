package main

import (
	"github.com/yuin/gopher-lua"
)

func runnerModeLoader(L *lua.LState) *lua.LTable {
	exports := map[string]lua.LGFunction{
		"sh": shRunner,
		"lua": luaRunner,
		"setMode": hlrunnerMode,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.SetField(mod, "mode", runnerMode)

	return mod
}

func shRunner(L *lua.LState) int {
	cmd := L.CheckString(1)
	exitCode, err := handleSh(cmd)
	var luaErr lua.LValue = lua.LNil
	if err != nil {
		luaErr = lua.LString(err.Error())
	}

	L.Push(lua.LNumber(exitCode))
	L.Push(luaErr)

	return 2
}

func luaRunner(L *lua.LState) int {
	cmd := L.CheckString(1)
	exitCode, err := handleLua(cmd)
	var luaErr lua.LValue = lua.LNil
	if err != nil {
		luaErr = lua.LString(err.Error())
	}

	L.Push(lua.LNumber(exitCode))
	L.Push(luaErr)

	return 2
}
