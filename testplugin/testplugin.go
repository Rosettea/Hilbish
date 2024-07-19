package main

import (
	rt "github.com/arnodel/golua/runtime"
)

func Loader(rtm *rt.Runtime) rt.Value {
	return rt.StringValue("hello world!")
}
