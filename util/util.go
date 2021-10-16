package util

import "github.com/yuin/gopher-lua"

func Document(L *lua.LState, module lua.LValue, doc string) {
	mt := L.NewTable()
	L.SetField(mt, "__doc", lua.LString(doc))

	L.SetMetatable(module, mt)
}
