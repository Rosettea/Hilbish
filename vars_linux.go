// +build linux

package main

// String vars that are free to be changed at compile time
var (
	requirePaths = `';./libs/?/?.lua;./libs/?/init.lua;./?/init.lua;./?/?.lua'
	.. ';/usr/share/hilbish/libs/?/init.lua;'
	.. ';/usr/share/hilbish/libs/?/?.lua;'
	.. hilbish.xdg.data .. '/hilbish/libs/?/init.lua;'
	.. hilbish.xdg.data	.. '/hilbish/libs/?/?.lua;'
	.. hilbish.xdg.data	.. '/hilbish/libs/?.lua'
	.. hilbish.xdg.config	.. '/hilbish/?/init.lua'
	.. hilbish.xdg.config	.. '/hilbish/?/?.lua'
	.. hilbish.xdg.config	.. '/hilbish/?.lua'`
	preloadPath = "/usr/share/hilbish/preload.lua"
	sampleConfPath = "/usr/share/hilbish/.hilbishrc.lua" // Path to default/sample config
)
