-- Default Hilbish config
local hilbish = require 'hilbish'
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

--[[
hilbish.timeout(function()
	hilbish.messages.send {title = 'greetings!', text = 'hello world :D'}
end, 2000)
]]--

bait.catch('hilbish.notification', function()
	hilbish.prompt(lunacolors.blue('• 1 new notification'), 'right')

	hilbish.timeout(function()
		hilbish.prompt('', 'right')
	end, 3000)
end)
