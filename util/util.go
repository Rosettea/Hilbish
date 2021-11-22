package util

import "github.com/yuin/gopher-lua"
import "fmt"

// Document adds a documentation string to a module.
// It is accessible via the __doc metatable.
func Document(L *lua.LState, module lua.LValue, doc string) {
	mt := L.GetMetatable(module)
	if mt == lua.LNil {
		mt = L.NewTable()
		docProp := L.NewTable()
		L.SetField(mt, "__docProp", docProp)

		L.SetMetatable(module, mt)
	}
	L.SetField(mt, "__doc", lua.LString(doc))
}

// SetField sets a field in a table, adding docs for it.
// It is accessible via the __docProp metatable. It is a table of the names of the fields.
func SetField(L *lua.LState, module lua.LValue, field string, value lua.LValue, doc string) {
	mt := L.GetMetatable(module)
	docProp := L.GetTable(mt, lua.LString("__docProp"))
	fmt.Println("docProp", docProp)

	L.SetField(docProp, field, lua.LString(doc))
	L.SetField(module, field, value)
}
