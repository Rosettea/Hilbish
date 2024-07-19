package moonlight

import (
	rt "github.com/arnodel/golua/runtime"
)

type Value = rt.Value

func StringValue(str string) Value {
	return rt.StringValue(str)
}

func IntValue(i int) Value {
	return rt.IntValue(int64(i))
}

func BoolValue(b bool) Value {
	return rt.BoolValue(b)
}

func TableValue(t *Table) Value {
	return rt.TableValue(t.lt)
}
