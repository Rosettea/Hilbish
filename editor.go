package main

import (
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

// #interface editor
// interactions for Hilbish's line reader
// The hilbish.editor interface provides functions to
// directly interact with the line editor in use.
func editorLoader(rtm *rt.Runtime) *rt.Table {
	exports := map[string]util.LuaExport{
		"insert": {editorInsert, 1, false},
		"setVimRegister": {editorSetRegister, 1, false},
		"getVimRegister": {editorGetRegister, 2, false},
		"getLine": {editorGetLine, 0, false},
	}

	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	return mod
}

// #interface editor
// insert(text)
// Inserts text into the line.
func editorInsert(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
/*
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	text, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	lr.rl.Insert(text)
*/

	return c.Next(), nil
}

// #interface editor
// setVimRegister(register, text)
// Sets the vim register at `register` to hold the passed text.
// --- @param register string
// --- @param text string
func editorSetRegister(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
/*
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
*/

	return c.Next(), nil
}

// #interface editor
// getVimRegister(register)
// Returns the text that is at the register.
// --- @param register string
func editorGetRegister(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	/*
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	register, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	buf := lr.rl.GetFromRegister(register)

	return c.PushingNext1(t.Runtime, rt.StringValue(string(buf))), nil
	*/

	return c.Next(), nil
}

// #interface editor
// getLine()
// Returns the current input line.
func editorGetLine(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	/*
	buf := lr.rl.GetLine()

	return c.PushingNext1(t.Runtime, rt.StringValue(string(buf))), nil
	*/
	return c.Next(), nil
}
