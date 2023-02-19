-- Greenhouse is a simple text scrolling handler for terminal programs.
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
	self.pages = {}
	self.curPage = 1

	return self
end

function Greenhouse:addPage(page)
	table.insert(self.pages, page)
end

function Greenhouse:updateCurrentPage(text)
	local page = self.pages[self.curPage]
	page:setText(text)
end

function Greenhouse:draw()
	local lines = self.pages[self.curPage].lines
	self.sink:write(ansikit.getCSI(self.start .. ';1', 'H'))
	self.sink:write(ansikit.getCSI(2, 'J'))

	-- the -2 negate is for the command and status line
	for i = self.offset, self.offset + (self.region.height - self.start) - 2 do
		if i > #lines then break end
		self.sink:writeln('\r' .. lines[i]:gsub('\t', '        '):sub(0, self.region.width - 2))
	end
	self.sink:write '\r'

	self.sink:write(ansikit.getCSI(self.region.height - self.start.. ';1', 'H'))
	self.sink:writeln(string.format('Page %d', self.curPage))
end

function Greenhouse:scroll(direction)
	local lines = self.pages[self.curPage].lines

	local oldOffset = self.offset
	if direction == 'down' then
		self.offset = math.min(self.offset + 1, #lines)
	elseif direction == 'up' then
		self.offset = math.max(self.offset - 1, 1)
	end

	if self.offset ~= oldOffset then self:draw() end
end

function Greenhouse:update()
	local size = terminal.size()
	self.region = size

	self:draw()
end

function Greenhouse:next()
	local oldCurrent = self.curPage
	self.curPage = math.min(self.curPage + 1, #self.pages)
	if self.curPage ~= oldCurrent then
		self.offset = 1
		self:draw()
	end
end

function Greenhouse:previous()
	local oldCurrent = self.curPage
	self.curPage = math.max(self.curPage - 1, 1)
	if self.curPage ~= oldCurrent then
		self.offset = 1
		self:draw()
	end
end

return Greenhouse
