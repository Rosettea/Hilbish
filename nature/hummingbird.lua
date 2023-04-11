local bait = require 'bait'
local commander = require 'commander'
local lunacolors = require 'lunacolors'

local M = {}
local counter = 0
M._messages = {}
M.icons = {
	INFO = '',
	SUCCESS = '',
	WARN = '',
	ERROR = ''
}

hilbish.messages = {}

--- Represents a Hilbish message.
--- @class hilbish.message
--- @field icon string Unicode (preferably standard emoji) icon for the message notification.
--- @field title string Title of the message (like an email subject).
--- @field text string Contents of the message.
--- @field channel string Short identifier of the message. `hilbish` and `hilbish.*` is preserved for internal Hilbish messages.
--- @field summary string A short summary of the message.
--- @field read boolean Whether the full message has been read or not.

function expect(tbl, field)
	if not tbl[field] or tbl[field] == '' then
		error(string.format('expected field %s in message'))
	end
end

--- Sends a message.
--- @param message hilbish.message
function hilbish.messages.send(message)
	expect(message, 'text')
	expect(message, 'title')
	counter = counter + 1
	message.index = counter

	M._messages[message.index] = message
	bait.throw('hilbish.notification', message)
end

function hilbish.messages.all()
	return M._messages
end

commander.register('messages', function(_, sinks)
	for _, msg in ipairs(hilbish.messages.all()) do
		local heading = lunacolors.format(string.format('Message {cyan}#%d{reset}: %s', msg.index, msg.title))
		sinks.out:writeln(heading)
		sinks.out:writeln(string.rep('=', string.len(heading)))
		sinks.out:writeln(msg.text)
	end
end)

return M
