package main

import (
	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
	"github.com/blackfireio/osinfo"
)

// #interface os
// OS Info
// The `os` interface provides simple text information properties about
// the current OS on the systen. This mainly includes the name and
// version.
// #property family Family name of the current OS
// #property name Pretty name of the current OS
// #property version Version of the current OS
func hshosLoader(rtm *rt.Runtime) *rt.Table {
	info, _ := osinfo.GetOSInfo()
	mod := rt.NewTable()

	util.SetField(rtm, mod, "family", rt.StringValue(info.Family), "Family name of the current OS")
	util.SetField(rtm, mod, "name", rt.StringValue(info.Name), "Pretty name of the current OS")
	util.SetField(rtm, mod, "version", rt.StringValue(info.Version), "Version of the current OS")

	return mod
}
