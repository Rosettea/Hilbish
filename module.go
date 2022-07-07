package main

import (
	"encoding/gob"
	"errors"
	"os"
	"os/exec"

	"hilbish/util"

	"github.com/Rosettea/Malvales"
	"github.com/hashicorp/go-plugin"
	rt "github.com/arnodel/golua/runtime"
)

func moduleLoader(rtm *rt.Runtime) *rt.Table {
	gob.Register(os.File{})
	exports := map[string]util.LuaExport{
		"load": {moduleLoad, 2, false},
	}

	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)

	return mod
}

func moduleLoad(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	
	path, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	name, err := c.StringArg(1)
	if err != nil {
		return nil, err
	}

	// plugin is just go executable; check if it is (or exists)
	if err := findExecutable(path, false, false); err != nil {
		return nil, err
	}

	moduleHandshake := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "HSH_PLUGIN",
		MagicCookieValue: name,
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: moduleHandshake,
		Plugins: map[string]plugin.Plugin{
			"entry": &malvales.Entry{},
		},
		Cmd: exec.Command(path),
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	ret, err := rpcClient.Dispense("entry")
	if err != nil {
		return nil, err
	}

	plug, ok := ret.(malvales.Plugin)
	if !ok {
		return nil, errors.New("did not get plugin from module")
	}

	val := plug.Loader(t.Runtime)

	return c.PushingNext1(t.Runtime, val), nil
}
