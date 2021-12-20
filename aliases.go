package main

import (
	"strings"
	"sync"

	"github.com/yuin/gopher-lua"
)

var aliases *hilbishAliases

type hilbishAliases struct {
	aliases map[string]string
	mu *sync.RWMutex
}

// initialize aliases map
func NewAliases() *hilbishAliases {
	return &hilbishAliases{
		aliases: make(map[string]string),
		mu: &sync.RWMutex{},
	}
}

func (h *hilbishAliases) Add(alias, cmd string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.aliases[alias] = cmd
}

func (h *hilbishAliases) All() map[string]string {
	return h.aliases
}

func (h *hilbishAliases) Delete(alias string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.aliases, alias)
}

func (h *hilbishAliases) Resolve(cmdstr string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	args := strings.Split(cmdstr, " ")
	for h.aliases[args[0]] != "" {
		alias := h.aliases[args[0]]
		cmdstr = alias + strings.TrimPrefix(cmdstr, args[0])
		cmdArgs, _ := splitInput(cmdstr)
		args = cmdArgs

		if h.aliases[args[0]] == alias {
			break
		}
		if h.aliases[args[0]] != "" {
			continue
		}
	}

	return cmdstr
}

// lua section

func (h *hilbishAliases) Loader(L *lua.LState) *lua.LTable {
	// create a lua module with our functions
	hshaliasesLua := map[string]lua.LGFunction{
		"add": h.luaAdd,
		"list": h.luaList,
		"del": h.luaDelete,
	}

	mod := L.SetFuncs(L.NewTable(), hshaliasesLua)

	return mod
}

func (h *hilbishAliases) luaAdd(L *lua.LState) int {
	alias := L.CheckString(1)
	cmd := L.CheckString(2)
	h.Add(alias, cmd)

	return 0
}

func (h *hilbishAliases) luaList(L *lua.LState) int {
	aliasesList := L.NewTable()
	for k, v := range h.All() {
		aliasesList.RawSetString(k, lua.LString(v))
	}

	L.Push(aliasesList)

	return 1
}

func (h *hilbishAliases) luaDelete(L *lua.LState) int {
	alias := L.CheckString(1)
	h.Delete(alias)

	return 0
}
