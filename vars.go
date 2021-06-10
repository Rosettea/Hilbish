package main

// String vars that are free to be changed at compile time
var (
	version = "v0.4.0"
	requirePaths = `';./libs/?/init.lua;./?/init.lua;./?/?.lua'
	.. ';/usr/share/hilbish/libs/?/init.lua;'
	.. ';/usr/share/hilbish/libs/?/?.lua;'
	.. os.getenv 'HOME' .. '/.local/share/hilbish/libs/?/init.lua;'
	.. os.getenv 'HOME'	.. '/.local/share/hilbish/libs/?/?.lua;'
	.. os.getenv 'HOME'	.. '/.local/share/hilbish/libs/?.lua'`
	preloadPath = "/usr/share/hilbish/preload.lua"
	defaultConfDir = "" // ~ will be substituted for home, path for user's default config
	sampleConfPath = "/usr/share/hilbish/.hilbishrc.lua" // Path to default/sample config

	prompt string // Prompt will always get changed anyway
	multilinePrompt = "> "
)

// Flags
var (
	running bool // Is a command currently running
	interactive bool
	login bool // Are we the login shell?
	noexecute bool // Should we run Lua or only report syntax errors
)

