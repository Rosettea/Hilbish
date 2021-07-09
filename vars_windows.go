// +build windows

package main

// String vars that are free to be changed at compile time
var (
	requirePaths = `';./libs/?/init.lua;./?/init.lua;./?/?.lua'
	.. hilbish.home .. '\\Appdata\\Roaming\\Hilbish\\libs\\?\\init.lua;'
	.. hilbish.home .. '\\Appdata\\Roaming\\Hilbish\\libs\\?\\?.lua;'`
	preloadPath = "~\\Appdata\\Roaming\\Hilbish\\preload.lua" // ~ and \ gonna cry?
	sampleConfPath = "~\\Appdata\\Roaming\\Hilbish\\hilbishrc.lua" // Path to default/sample config
)
