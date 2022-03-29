package util

import (
	"os"

	"github.com/yuin/gopher-lua"
	rt "github.com/arnodel/golua/runtime"
)

// Document adds a documentation string to a module.
// It is accessible via the __doc metatable.
func Document(L *lua.LState, module lua.LValue, doc string) {
/*
	mt := L.GetMetatable(module)
	if mt == lua.LNil {
		mt = L.NewTable()
		L.SetMetatable(module, mt)
	}
	L.SetField(mt, "__doc", lua.LString(doc))
*/
}

// SetField sets a field in a table, adding docs for it.
// It is accessible via the __docProp metatable. It is a table of the names of the fields.
func SetField(rtm *rt.Runtime, module *rt.Table, field string, value rt.Value, doc string) {
	// TODO:    ^ rtm isnt needed, i should remove it
	mt := module.Metatable()
	
	if mt == nil {
		mt = rt.NewTable()
		docProp := rt.NewTable()
		mt.Set(rt.StringValue("__docProp"), rt.TableValue(docProp))

		module.SetMetatable(mt)
	}
	docProp := mt.Get(rt.StringValue("__docProp"))

	docProp.AsTable().Set(rt.StringValue(field), rt.StringValue(doc))
	module.Set(rt.StringValue(field), value)
}

func DoString(rtm *rt.Runtime, code string) error {
	chunk, err := rtm.CompileAndLoadLuaChunk("", []byte(code), rt.TableValue(rtm.GlobalEnv()))
	if chunk != nil {
		_, err = rt.Call1(rtm.MainThread(), rt.FunctionValue(chunk))
	}

	return err
}

func DoFile(rtm *rt.Runtime, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return DoString(rtm, string(data))
}

func HandleStrCallback(t *rt.Thread, c *rt.GoCont) (string, *rt.Closure, error) {
	if err := c.CheckNArgs(2); err != nil {
		return "", nil, err
	}
	name, err := c.StringArg(0)
	if err != nil {
		return "", nil, err
	}
	cb, err := c.ClosureArg(1)
	if err != nil {
		return "", nil, err
	}

	return name, cb, err
}
