-- Greenhouse is a simple text scrolling handler for terminal program.
-- The idea is that it can be set a region to do its scrolling and paging
-- job and then the user can draw whatever outside it.
-- This reduces code duplication for the message viewer
-- and flowerbook.

local ansikit = require 'ansikit'
local terminal = require 'terminal'
local Object = require 'nature.object'

local Greenhouse = Object:extend()

function Greenhouse:new(sink)
	local size = terminal.size()
	self.region = size
	self.start = 1
	self.offset = 1
	self.sink = sink

	return self
end

function Greenhouse:setText(text)
	self.lines = string.split(text, '\n')
end

function Greenhouse:draw()
	self.sink:write(ansikit.getCSI(self.start .. ';1', 'H'))

	for i = self.offset, self.offset + (self.region.height - self.start) - 1 do
		self.sink:write(ansikit.getCSI(2, 'K'))
		self.sink:writeln(self.lines[i]:gsub('\t', '        '):sub(0, self.region.width - 2))
	end
end

function Greenhouse:scroll(direction)
	if direction == 'down' then
		self.offset = math.min(self.offset + 1, #self.lines)
	elseif direction == 'up' then
		self.offset = math.max(self.offset - 1, 1)
	end
	self:draw()
end

return Greenhouse
