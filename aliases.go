package main

import (
	"strings"
	"sync"

	"github.com/yuin/gopher-lua"
)

var aliases *aliasHandler

type aliasHandler struct {
	aliases map[string]string
	mu *sync.RWMutex
}

// initialize aliases map
func newAliases() *aliasHandler {
	return &aliasHandler{
		aliases: make(map[string]string),
		mu: &sync.RWMutex{},
	}
}

func (a *aliasHandler) Add(alias, cmd string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.aliases[alias] = cmd
}

func (a *aliasHandler) All() map[string]string {
	return a.aliases
}

func (a *aliasHandler) Delete(alias string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.aliases, alias)
}

func (a *aliasHandler) Resolve(cmdstr string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	args := strings.Split(cmdstr, " ")
	for a.aliases[args[0]] != "" {
		alias := a.aliases[args[0]]
		cmdstr = alias + strings.TrimPrefix(cmdstr, args[0])
		cmdArgs, _ := splitInput(cmdstr)
		args = cmdArgs

		if a.aliases[args[0]] == alias {
			break
		}
		if a.aliases[args[0]] != "" {
			continue
		}
	}

	return cmdstr
}

// lua section

func (a *aliasHandler) Loader(L *lua.LState) *lua.LTable {
	// create a lua module with our functions
	hshaliasesLua := map[string]lua.LGFunction{
		"add": a.luaAdd,
		"list": a.luaList,
		"del": a.luaDelete,
	}

	mod := L.SetFuncs(L.NewTable(), hshaliasesLua)

	return mod
}

func (a *aliasHandler) luaAdd(L *lua.LState) int {
	alias := L.CheckString(1)
	cmd := L.CheckString(2)
	a.Add(alias, cmd)

	return 0
}

func (a *aliasHandler) luaList(L *lua.LState) int {
	aliasesList := L.NewTable()
	for k, v := range a.All() {
		aliasesList.RawSetString(k, lua.LString(v))
	}

	L.Push(aliasesList)

	return 1
}

func (a *aliasHandler) luaDelete(L *lua.LState) int {
	alias := L.CheckString(1)
	a.Delete(alias)

	return 0
}
