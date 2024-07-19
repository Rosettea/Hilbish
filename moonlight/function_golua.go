package moonlight

import (
	rt "github.com/arnodel/golua/runtime"
)

type GoFunctionFunc = rt.GoFunctionFunc

type GoCont = rt.GoCont
type Cont = rt.Cont

func (mlr *Runtime) CheckNArgs(c *GoCont, num int) error {
	return c.CheckNArgs(num)
}

func (mlr *Runtime) StringArg(c *GoCont, num int) (string, error) {
	return c.StringArg(num)
}

func (mlr *Runtime) GoFunction(fun GoToLuaFunc) rt.GoFunctionFunc {
	return func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		gocont := GoCont(*c)
		return fun(mlr, &gocont)
	}
}

func (mlr *Runtime) Call1(val Value, args ...Value) (Value, error) {
	return rt.Call1(mlr.rt.MainThread(), val, args...)
}
