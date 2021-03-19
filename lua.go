package main

import (
	"github.com/yuin/gopher-lua"
)

func hshprompt(L *lua.LState) int {
	prompt = L.ToString(1)

	return 0
}
