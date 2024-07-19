local ansikit = require 'ansikit'
local commander = require 'commander'

commander.register('clear', function()
	ansikit.clear(true)
	ansikit.cursorTo(0, 0)
end)
