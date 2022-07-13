// +build windows

package main

import "hilbish/util"

// String vars that are free to be changed at compile time
var (
	requirePaths = commonRequirePaths + `.. ';'
	.. hilbish.userDir.config .. '\\Hilbish\\libs\\?\\init.lua;'
	.. hilbish.userDir.config .. '\\Hilbish\\libs\\?\\?.lua;'
	.. hilbish.userDir.config .. '\\Hilbish\\libs\\?.lua;'`
	dataDir = util.ExpandHome("~\\Appdata\\Roaming\\Hilbish") // ~ and \ gonna cry?
	preloadPath = dataDir + "\\nature\\init.lua"
	sampleConfPath = dataDir + "\\hilbishrc.lua" // Path to default/sample config
	defaultConfDir = ""
)
