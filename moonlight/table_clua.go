//go:build midnight
package moonlight

import (
	"fmt"

	"github.com/aarzilli/golua/lua"
)

type Table struct{
	refIdx int
	mlr *Runtime
	nativeFields map[Value]Value
}

func NewTable() *Table {
	return &Table{
		refIdx: -1,
		nativeFields: make(map[Value]Value),
	}
}

func (t *Table) SetRuntime(mlr *Runtime) {
	t.mlr = mlr

	if t.refIdx == -1 {
		mlr.state.NewTable()
		t.refIdx = mlr.state.Ref(lua.LUA_REGISTRYINDEX)
		mlr.state.Pop(1)
	}
}

func (t *Table) Get(val Value) Value {
	return NilValue
}

func (t *Table) Push() {
	t.mlr.state.RawGeti(lua.LUA_REGISTRYINDEX, t.refIdx)
}

func (t *Table) SetField(key string, value Value) {
	fmt.Printf("key: %s, value: %s\n", key, value.TypeName())
	t.Push()
	defer t.mlr.state.Pop(1)

	t.mlr.pushToState(value)
	t.mlr.state.SetField(-1, key)
	t.mlr.state.Pop(1)
	println("what")
}

func (t *Table) Set(key Value, value Value) {
	t.nativeFields[key] = value
}

func ForEach(tbl *Table, cb func(key Value, val Value)) {
}

func (mlr *Runtime) GlobalTable() *Table {
	mlr.state.GetGlobal("_G")
	return &Table{
		refIdx: mlr.state.Ref(lua.LUA_REGISTRYINDEX),
	}
}

func ToTable(v Value) *Table {
	return &Table{
		refIdx: -1,
	}
}

func TryTable(v Value) (*Table, bool) {
	return nil, false
}

func (t *Table) setRefIdx(mlr *Runtime, i int) {
	t.refIdx = mlr.state.Ref(i)
}
