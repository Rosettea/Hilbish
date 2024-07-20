package moonlight

import (
	rt "github.com/arnodel/golua/runtime"
)

type Table struct{
	lt *rt.Table
}

func NewTable() *Table {
	return &Table{
		lt: rt.NewTable(),
	}
}

func (t *Table) Get(val Value) Value {
	return t.lt.Get(val)
}

func (t *Table) SetField(key string, value Value) {
	t.lt.Set(rt.StringValue(key), value)
}

func (t *Table) Set(key Value, value Value) {
	t.lt.Set(key, value)
}

func (mlr *Runtime) GlobalTable() *Table {
	return &Table{
		lt: mlr.rt.GlobalEnv(),
	}
}
