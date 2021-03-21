package main

import (
	"github.com/yuin/gopher-lua"
)

func hshprompt(L *lua.LState) int {
	prompt = L.ToString(1)

	return 0
}

func hshalias(L *lua.LState) int {
	alias := L.ToString(1)
	source := L.ToString(2)

	aliases[alias] = source

	return 1
}
