package main

import (
	"bufio"
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
	writer *bufio.Writer
	reader *bufio.Reader
	ud *rt.UserData
	autoFlush bool
}

func setupSinkType(rtm *rt.Runtime) {
	sinkMeta := rt.NewTable()

	sinkMethods := rt.NewTable()
	sinkFuncs := map[string]util.LuaExport{
		"read": {luaSinkRead, 0, false},
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
// read() -> string
// --- @returns string
// Reads input from the sink.
func luaSinkRead(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	s, err := sinkArg(c, 0)
	if err != nil {
		return nil, err
	}

	str, _ := s.reader.ReadString('\n')

	return c.PushingNext1(t.Runtime, rt.StringValue(str)), nil
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
	if s.autoFlush {
		s.writer.Flush()
	}

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
	if s.autoFlush {
		s.writer.Flush()
	}

	return c.Next(), nil
}

// #member
// flush()
// Flush writes all buffered input to the sink.
func luaSinkFlush(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}

	s, err := sinkArg(c, 0)
	if err != nil {
		return nil, err
	}

	s.writer.Flush()

	return c.Next(), nil
}

func newSinkInput(r io.Reader) *sink {
	s := &sink{
		reader: bufio.NewReader(r),
	}
	s.ud = sinkUserData(s)

	return s
}

func newSinkOutput(w io.Writer) *sink {
	s := &sink{
		writer: bufio.NewWriter(w),
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
