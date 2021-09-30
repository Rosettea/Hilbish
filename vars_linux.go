// +build linux

package main

// String vars that are free to be changed at compile time
var (
	requirePaths = `';./libs/?/?.lua;./libs/?/init.lua;./?/init.lua;./?/?.lua'
	.. ';/usr/share/hilbish/libs/?/init.lua;'
	.. ';/usr/share/hilbish/libs/?/?.lua;'
	.. hilbish.home .. '/.local/share/hilbish/libs/?/init.lua;'
	.. hilbish.home	.. '/.local/share/hilbish/libs/?/?.lua;'
	.. hilbish.home	.. '/.local/share/hilbish/libs/?.lua'
	.. hilbish.home	.. '/.config/hilbish/?/init.lua'
	.. hilbish.home	.. '/.config/hilbish/?/?.lua'
	.. hilbish.home	.. '/.config/hilbish/?.lua'`
	preloadPath = "/usr/share/hilbish/preload.lua"
	sampleConfPath = "/usr/share/hilbish/.hilbishrc.lua" // Path to default/sample config
)
