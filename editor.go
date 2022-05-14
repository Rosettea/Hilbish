package main

import (
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

func editorLoader(rtm *rt.Runtime) *rt.Table {
	exports := map[string]util.LuaExport{
		"insert": {editorInsert, 1, false},
		"setVimRegister": {editorSetRegister, 1, false},
		"getVimRegister": {editorGetRegister, 2, false},
	}

	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	return mod
}

func editorInsert(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	text, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	lr.rl.Insert(text)

	return c.Next(), nil
}

func editorSetRegister(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	register, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	text, err := c.StringArg(1)
	if err != nil {
		return nil, err
	}

	lr.rl.SetRegisterBuf(register, []rune(text))

	return c.Next(), nil
}

func editorGetRegister(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	register, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	buf := lr.rl.GetFromRegister(register)

	return c.PushingNext1(t.Runtime, rt.StringValue(string(buf))), nil
}
