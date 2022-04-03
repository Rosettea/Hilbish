package util

import (
	"bufio"
	"io"
	"os"

	rt "github.com/arnodel/golua/runtime"
)

// Document adds a documentation string to a module.
// It is accessible via the __doc metatable.
func Document(module *rt.Table, doc string) {
	mt := module.Metatable()
	
	if mt == nil {
		mt = rt.NewTable()
		module.SetMetatable(mt)
	}

	mt.Set(rt.StringValue("__doc"), rt.StringValue(doc))
}

// SetField sets a field in a table, adding docs for it.
// It is accessible via the __docProp metatable. It is a table of the names of the fields.
func SetField(rtm *rt.Runtime, module *rt.Table, field string, value rt.Value, doc string) {
	// TODO:    ^ rtm isnt needed, i should remove it
	mt := module.Metatable()
	
	if mt == nil {
		mt = rt.NewTable()
		docProp := rt.NewTable()
		mt.Set(rt.StringValue("__docProp"), rt.TableValue(docProp))

		module.SetMetatable(mt)
	}
	docProp := mt.Get(rt.StringValue("__docProp"))

	docProp.AsTable().Set(rt.StringValue(field), rt.StringValue(doc))
	module.Set(rt.StringValue(field), value)
}

// DoString runs the code string in the Lua runtime.
func DoString(rtm *rt.Runtime, code string) error {
	chunk, err := rtm.CompileAndLoadLuaChunk("<string>", []byte(code), rt.TableValue(rtm.GlobalEnv()))
	if chunk != nil {
		_, err = rt.Call1(rtm.MainThread(), rt.FunctionValue(chunk))
	}

	return err
}

// DoFile runs the contents of the file in the Lua runtime.
func DoFile(rtm *rt.Runtime, path string) error {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return err
	}
	
	reader := bufio.NewReader(f)
	c, err := reader.ReadByte()
	if err != nil && err != io.EOF {
		return err
	}

	// unread so a char won't be missing
	err = reader.UnreadByte()
	if err != nil {
		return err
	}

	var buf []byte
	if c == byte('#') {
		// shebang - skip that line
		_, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		buf = []byte{'\n'}
	}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		
		buf = append(buf, line...)
	}

	chunk, err := rtm.CompileAndLoadLuaChunk(path, buf, rt.TableValue(rtm.GlobalEnv()))
	if chunk != nil {
		_, err = rt.Call1(rtm.MainThread(), rt.FunctionValue(chunk))
	}

	return err
}

// HandleStrCallback handles function parameters for Go functions which take
// a string and a closure.
func HandleStrCallback(t *rt.Thread, c *rt.GoCont) (string, *rt.Closure, error) {
	if err := c.CheckNArgs(2); err != nil {
		return "", nil, err
	}
	name, err := c.StringArg(0)
	if err != nil {
		return "", nil, err
	}
	cb, err := c.ClosureArg(1)
	if err != nil {
		return "", nil, err
	}

	return name, cb, err
}
