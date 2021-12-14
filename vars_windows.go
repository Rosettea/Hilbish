// +build windows

package main

// String vars that are free to be changed at compile time
var (
	requirePaths = commonRequirePaths + `.. ';'
	.. hilbish.user.config .. '\\Hilbish\\libs\\?\\init.lua;'
	.. hilbish.user.config .. '\\Hilbish\\libs\\?\\?.lua;'
	.. hilbish.user.config .. '\\Hilbish\\libs\\?.lua;'`
	dataDir = "~\\Appdata\\Roaming\\Hilbish" // ~ and \ gonna cry?
	preloadPath = dataDir + "\\preload.lua"
	sampleConfPath = dataDir + "\\hilbishrc.lua" // Path to default/sample config
)
