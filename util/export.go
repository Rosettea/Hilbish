package util

import (
	rt "github.com/arnodel/golua/runtime"
)

type LuaExport struct {
	Function rt.GoFunctionFunc
	ArgNum int
	Variadic bool
}

func SetExports(rtm *rt.Runtime, tbl *rt.Table, exports map[string]LuaExport) {
	for name, export := range exports {
		rtm.SetEnvGoFunc(tbl, name, export.Function, export.ArgNum, export.Variadic)
	}
}
