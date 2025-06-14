//go:build !midnight

package moonlight

import (
	rt "github.com/arnodel/golua/runtime"
)

type GoFunctionFunc = rt.GoFunctionFunc

func (mlr *Runtime) CheckNArgs(num int) error {
	return mlr.rt.MainThread().CurrentCont().(*rt.GoCont).CheckNArgs(num)
}

func (mlr *Runtime) Check1Arg() error {
	return mlr.rt.MainThread().CurrentCont().(*rt.GoCont).Check1Arg()
}

func (mlr *Runtime) StringArg(num int) (string, error) {
	return mlr.rt.MainThread().CurrentCont().(*rt.GoCont).StringArg(num)
}

func (mlr *Runtime) ClosureArg(num int) (*Closure, error) {
	return mlr.rt.MainThread().CurrentCont().(*rt.GoCont).ClosureArg(num)
}

func (mlr *Runtime) Arg(c *GoCont, num int) Value {
	return mlr.rt.MainThread().CurrentCont().(*rt.GoCont).Arg(num)
}

func (mlr *Runtime) GoFunction(fun GoToLuaFunc) GoFunctionFunc {
	return func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		return c.Next(), fun(mlr)
	}
}
