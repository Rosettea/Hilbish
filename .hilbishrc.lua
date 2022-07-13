-- Default Hilbish config
local lunacolors = require 'lunacolors'
local bait = require 'bait'
local ansikit = require 'ansikit'

local function doPrompt(fail)
	hilbish.prompt(lunacolors.format(
		'{blue}%u {cyan}%d ' .. (fail and '{red}' or '{green}') .. '∆ '
	))
end

doPrompt()

bait.catch('command.exit', function(code)
	doPrompt(code ~= 0)
end)

bait.catch('hilbish.vimMode', function(mode)
	if mode ~= 'insert' then
		ansikit.cursorStyle(ansikit.blockCursor)
	else
		ansikit.cursorStyle(ansikit.lineCursor)
	end
end)
