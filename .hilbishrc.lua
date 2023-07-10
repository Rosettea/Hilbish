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

bait.catch('hilbish.notification', function()
	local notif = string.format('• %s unread notification%s', hilbish.messages.unreadCount(), hilbish.messages.unreadCount() > 1 and 's' or '')
	hilbish.prompt(lunacolors.blue(notif), 'right')

	hilbish.timeout(function()
		hilbish.prompt('', 'right')
	end, 3000)
end)
