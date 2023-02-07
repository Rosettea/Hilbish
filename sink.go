package main

import (
	"fmt"
	"io"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
)

var sinkMetaKey = rt.StringValue("hshsink")

// #type
// A sink is a structure that has input and/or output to/from
// a desination.
type sink struct{
	writer io.Writer
	reader io.Reader
	ud *rt.UserData
}

func setupSinkType(rtm *rt.Runtime) {
	sinkMeta := rt.NewTable()

	sinkMethods := rt.NewTable()
	sinkFuncs := map[string]util.LuaExport{
		"write": {luaSinkWrite, 2, false},
		"writeln": {luaSinkWriteln, 2, false},
	}
	util.SetExports(l, sinkMethods, sinkFuncs)

	sinkIndex := func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		arg := c.Arg(1)
		val := sinkMethods.Get(arg)

		return c.PushingNext1(t.Runtime, val), nil
	}

	sinkMeta.Set(rt.StringValue("__index"), rt.FunctionValue(rt.NewGoFunction(sinkIndex, "__index", 2, false)))
	l.SetRegistry(sinkMetaKey, rt.TableValue(sinkMeta))
}

// #member
// write(str)
// Writes data to a sink.
func luaSinkWrite(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}

	s, err := sinkArg(c, 0)
	if err != nil {
		return nil, err
	}
	data, err := c.StringArg(1)
	if err != nil {
		return nil, err
	}

	s.writer.Write([]byte(data))

	return c.Next(), nil
}

// #member
// writeln(str)
// Writes data to a sink with a newline at the end.
func luaSinkWriteln(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}

	s, err := sinkArg(c, 0)
	if err != nil {
		return nil, err
	}
	data, err := c.StringArg(1)
	if err != nil {
		return nil, err
	}

	s.writer.Write([]byte(data + "\n"))

	return c.Next(), nil
}

func newSinkInput(r io.Reader) *sink {
	s := &sink{
		reader: r,
	}
	s.ud = sinkUserData(s)

	return s
}

func newSinkOutput(w io.Writer) *sink {
	s := &sink{
		writer: w,
	}
	s.ud = sinkUserData(s)

	return s
}

func sinkArg(c *rt.GoCont, arg int) (*sink, error) {
	s, ok := valueToSink(c.Arg(arg))
	if !ok {
		return nil, fmt.Errorf("#%d must be a sink", arg + 1)
	}

	return s, nil
}

func valueToSink(val rt.Value) (*sink, bool) {
	u, ok := val.TryUserData()
	if !ok {
		return nil, false
	}

	s, ok := u.Value().(*sink)
	return s, ok
}

func sinkUserData(s *sink) *rt.UserData {
	sinkMeta := l.Registry(sinkMetaKey)
	return rt.NewUserData(s, sinkMeta.AsTable())
}
