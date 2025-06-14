//go:build midnight

package moonlight

import (
	"fmt"

	"github.com/aarzilli/golua/lua"
)

type GoFunctionFunc struct {
	cf lua.LuaGoFunction
}

func (gf GoFunctionFunc) Continuation(mlr *Runtime, c Cont) Cont {
	return &GoCont{
		f:    gf,
		vals: []Value{},
	}
}

func (mlr *Runtime) CheckNArgs(num int) error {
	args := mlr.state.GetTop()
	if args < num {
		return fmt.Errorf("%d arguments needed", num)
	}

	return nil
}

func (mlr *Runtime) Check1Arg() error {
	return mlr.CheckNArgs(1)
}

func (mlr *Runtime) StringArg(num int) (string, error) {
	return mlr.state.CheckString(num + 1), nil
}

func (mlr *Runtime) Arg(c *GoCont, num int) Value {
	return c.vals[num]
}

func (mlr *Runtime) GoFunction(fun GoToLuaFunc) *GoFunctionFunc {
	mlr.returnNum = 0

	return &GoFunctionFunc{
		cf: func(L *lua.State) int {
			err := fun(mlr)
			if err != nil {
				L.RaiseError(err.Error())
				return 0
			}

			/*for _, val := range cont.(*GoCont).vals {
				switch Type(val) {
				case StringType:
					L.PushString(val.AsString())
				}
			}*/

			//return len(cont.(*GoCont).vals)
			return mlr.returnNum
		},
	}
}
