package main

import (
	"github.com/Rosettea/Malvales"
	"github.com/hashicorp/go-plugin"
	rt "github.com/arnodel/golua/runtime"
)

type TestPlugin struct {}

func (t *TestPlugin) Loader(rtm *rt.Runtime) rt.Value {
	println("hello")
	return rt.StringValue("hello world!")
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "HSH_PLUGIN",
	MagicCookieValue: "testplugin",
}

func main() {
	test := &TestPlugin{}

	var pluginMap = map[string]plugin.Plugin{
		"entry": &malvales.Entry{P: test},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
