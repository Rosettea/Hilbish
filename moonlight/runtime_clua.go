//go:build midnight

package moonlight

import (
	"fmt"

	"github.com/aarzilli/golua/lua"
)

type Runtime struct {
	state     *lua.State
	returnNum int
}

func NewRuntime() *Runtime {
	L := lua.NewState()
	L.OpenLibs()

	return &Runtime{
		state: L,
	}
}

func (mlr *Runtime) PushNext1(v Value) {
	mlr.returnNum = 1

	mlr.pushToState(v)
}

func (mlr *Runtime) Call1(f Value, args ...Value) (Value, error) {
	for _, arg := range args {
		mlr.pushToState(arg)
	}

	if f.refIdx > 0 {
		mlr.state.RawGeti(lua.LUA_REGISTRYINDEX, f.refIdx)
		mlr.state.Call(len(args), 1)
	}

	if mlr.state.GetTop() == 0 {
		return NilValue, nil
	}

	return NilValue, nil
}

func (mlr *Runtime) pushToState(v Value) {
	switch v.Type() {
	case NilType:
		mlr.state.PushNil()
	case StringType:
		mlr.state.PushString(v.AsString())
	case IntType:
		mlr.state.PushInteger(v.AsInt())
	case BoolType:
		mlr.state.PushBoolean(v.AsBool())
	case TableType:
		tbl := v.AsTable()
		tbl.SetRuntime(mlr)
		tbl.Push()
	case FunctionType:
		mlr.state.PushGoClosure(v.AsLuaFunction())
	default:
		fmt.Println("PUSHING UNIMPLEMENTED TYPE", v.TypeName())
		mlr.state.PushNil()
	}
}
