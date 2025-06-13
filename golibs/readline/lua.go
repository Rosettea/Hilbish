// line reader library
// The readline module is responsible for reading input from the user.
// The readline module is what Hilbish uses to read input from the user,
// including all the interactive features of Hilbish like history search,
// syntax highlighting, everything. The global Hilbish readline instance
// is usable at `hilbish.editor`.
package readline

import (
	"fmt"
	"io"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

var rlMetaKey = rt.StringValue("__readline")

func (rl *Instance) luaLoader(rtm *rt.Runtime) (rt.Value, func()) {
	rlMethods := rt.NewTable()
	rlMethodss := map[string]util.LuaExport{
		"read":   {luaRead, 1, false},
		"insert": {luaInsert, 2, false},
	}
	util.SetExports(rtm, rlMethods, rlMethodss)

	jobMeta := rt.NewTable()
	rlIndex := func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		_, err := rlArg(c, 0)
		if err != nil {
			return nil, err
		}

		arg := c.Arg(1)
		val := rlMethods.Get(arg)

		return c.PushingNext1(t.Runtime, val), nil
	}

	jobMeta.Set(rt.StringValue("__index"), rt.FunctionValue(rt.NewGoFunction(rlIndex, "__index", 2, false)))
	rtm.SetRegistry(rlMetaKey, rt.TableValue(jobMeta))

	rlFuncs := map[string]util.LuaExport{
		"new": {luaNew, 0, false},
	}

	luaRl := rt.NewTable()
	util.SetExports(rtm, luaRl, rlFuncs)

	return rt.TableValue(luaRl), nil
}

func luaNew(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	rl := NewInstance()
	ud := rlUserData(t.Runtime, rl)

	return c.PushingNext1(t.Runtime, rt.UserDataValue(ud)), nil
}

func luaInsert(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}

	rl, err := rlArg(c, 0)
	if err != nil {
		return nil, err
	}

	text, err := c.StringArg(1)
	if err != nil {
		return nil, err
	}

	rl.insert([]rune(text))

	return c.Next(), nil
}

func luaRead(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	rl, err := rlArg(c, 0)
	if err != nil {
		return nil, err
	}

	inp, err := rl.Readline()
	if err == EOF {
		fmt.Println("")
		return nil, io.EOF
	} else if err != nil {
		return nil, err
	}

	return c.PushingNext1(t.Runtime, rt.StringValue(inp)), nil
}

func rlArg(c *rt.GoCont, arg int) (*Instance, error) {
	j, ok := valueToRl(c.Arg(arg))
	if !ok {
		return nil, fmt.Errorf("#%d must be a readline", arg+1)
	}

	return j, nil
}

func valueToRl(val rt.Value) (*Instance, bool) {
	u, ok := val.TryUserData()
	if !ok {
		return nil, false
	}

	j, ok := u.Value().(*Instance)
	return j, ok
}

func rlUserData(rtm *rt.Runtime, rl *Instance) *rt.UserData {
	rlMeta := rtm.Registry(rlMetaKey)
	return rt.NewUserData(rl, rlMeta.AsTable())
}
