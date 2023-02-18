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
	self.offset = 0
	self.sink = sink

	return self
end

function Greenhouse:setText(text)
	self.lines = string.split(text, '\n')
end

function Greenhouse:draw()
	self.sink:write(ansikit.getCSI(self.start .. ';1', 'H'))

	for i = 1, #self.lines do
		if i > self.region.height - 1 then break end
		if not self.lines[i + self.offset] then break end

		self.sink:writeln(self.lines[i + self.offset]:gsub('\t', '        '):sub(0, self.region.width - 2))
	end
end

function Greenhouse:scroll(direction)
	if direction == 'down' then
		self.offset = self.offset + 1
	elseif direction == 'up' then
		self.offset = self.offset - 1
	end
	self:draw()
end

return Greenhouse
