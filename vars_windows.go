// +build windows

package main

// String vars that are free to be changed at compile time
var (
	requirePaths = `';./libs/?/init.lua;./?/init.lua;./?/?.lua'
	.. hilbish.home .. '\\Appdata\\Roaming\\Hilbish\\libs\\?\\init.lua;'
	.. hilbish.home .. '\\Appdata\\Roaming\\Hilbish\\libs\\?\\?.lua;'`
	dataDir = "~\\Appdata\\Roaming\\Hilbish" // ~ and \ gonna cry?
	preloadPath = dataDir + "\\preload.lua"
	sampleConfPath = dataDir + "\\hilbishrc.lua" // Path to default/sample config
)
