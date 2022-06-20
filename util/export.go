package util

import (
	rt "github.com/arnodel/golua/runtime"
)

// Interface is a Hilbish API interface.
type Interface struct{
	Name string
	Description string
	Setup func() *rt.Table
}

type Interfacer struct{
	ifaces map[string]*Interface
	loaded []string
	mod *rt.Table
}

func NewInterfacer(mod *rt.Table) *Interfacer {
	return &Interfacer{
		ifaces: make(map[string]*Interface),
		loaded: []string{},
		mod: mod,
	}
}

func (i *Interfacer) Add(ifaces []*Interface) {
	for _, iface := range ifaces {
		i.ifaces[iface.Name] = iface
	}
	// "claim" the name in interfacer module, to make sure user cant
	// overrides (like in hilbish table)
//	i.mod.Set(rt.StringValue(iface.Name), rt.BoolValue(true))
}

func (i *Interfacer) Load(name string) {
	if iface := i.ifaces[name]; iface != nil && !contains(i.loaded, name) {
		println(name)
		mod := iface.Setup()
		Document(mod, iface.Description)
		i.mod.Set(rt.StringValue(iface.Name), rt.TableValue(mod))

		i.loaded = append(i.loaded, iface.Name)
	}
}

// LuaExport represents a Go function which can be exported to Lua.
type LuaExport struct {
	Function rt.GoFunctionFunc
	ArgNum int
	Variadic bool
}

// SetExports puts the Lua function exports in the table.
func SetExports(rtm *rt.Runtime, tbl *rt.Table, exports map[string]LuaExport) {
	for name, export := range exports {
		rtm.SetEnvGoFunc(tbl, name, export.Function, export.ArgNum, export.Variadic)
	}
}
