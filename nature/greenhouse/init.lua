-- Greenhouse is a simple text scrolling handler for terminal programs.
-- The idea is that it can be set a region to do its scrolling and paging
-- job and then the user can draw whatever outside it.
-- This reduces code duplication for the message viewer
-- and flowerbook.

local ansikit = require 'ansikit'
local lunacolors = require 'lunacolors'
local terminal = require 'terminal'
local Page = require 'nature.greenhouse.page'
local Object = require 'nature.object'

local Greenhouse = Object:extend()

function Greenhouse:new(sink)
	local size = terminal.size()
	self.region = size
	self.contents = nil -- or can be a table
	self.start = 1 -- where to start drawing from (should replace with self.region.y)
	self.offset = 1 -- vertical text offset
	self.horizOffset = 1
	self.sink = sink
	self.pages = {}
	self.curPage = 1
	self.step = {
		horizontal = 5,
		vertical = 1
	}
	self.separator = '─'
	self.keybinds = {
		['Up'] = function(self) self:scroll 'up' end,
		['Down'] = function(self) self:scroll 'down' end,
		['Left'] = function(self) self:scroll 'left' end,
		['Right'] = function(self) self:scroll 'right' end,
		['Ctrl-Left'] = self.previous,
		['Ctrl-Right'] = self.next,
		['Ctrl-N'] = function(self) self:toc(true) end,
		['Enter'] = function(self)
			if self.isSpecial then
				self:jump(self.specialPageIdx)
				self:special(false)
			end
		end
	}
	self.isSpecial = false
	self.specialPage = nil
	self.specialPageIdx = 1
	self.specialOffset = 1

	return self
end

function Greenhouse:addPage(page)
	table.insert(self.pages, page)
end

function Greenhouse:updateCurrentPage(text)
	local page = self.pages[self.curPage]
	page:setText(text)
end

local function sub(str, offset, limit)
	local overhead = 0
	local function addOverhead(s)
		overhead = overhead + string.len(s)
	end

	local s = str:gsub('\x1b%[%d+;%d+;%d+;%d+;%d+%w', addOverhead)
     :gsub('\x1b%[%d+;%d+;%d+;%d+%w', addOverhead)
     :gsub('\x1b%[%d+;%d+;%d+%w',addOverhead)
     :gsub('\x1b%[%d+;%d+%w', addOverhead)
     :gsub('\x1b%[%d+%w', addOverhead)

	return s:sub(offset, utf8.offset(str, limit + overhead) or limit + overhead)
	--return s:sub(offset, limit + overhead)
end

function Greenhouse:draw()
	local workingPage = self.pages[self.curPage]
	local offset = self.offset
	if self.isSpecial then
		offset = self.specialOffset
		workingPage = self.specialPage
	end

	if workingPage.lazy and not workingPage.loaded then
		workingPage.initialize()
	end

	local lines = workingPage.lines
	self.sink:write(ansikit.getCSI(self.start .. ';1', 'H'))
	self.sink:write(ansikit.getCSI(2, 'J'))

	local writer = self.sink.writeln
	for i = offset, offset + self.region.height - 1 do
		if i > #lines then break end

		if i == offset + self.region.height - 1 then writer = self.sink.write end

		local line = lines[i]:gsub('{separator}', function() return self.separator:rep(self.region.width - 1) end)
		writer(self.sink, sub(line:gsub('\t', '        '), self.horizOffset, self.region.width))
	end
	writer(self.sink, '\27[0m')
	self:render()
end

function Greenhouse:render()
end

function Greenhouse:scroll(direction)
	if self.isSpecial then
		if direction == 'down' then
			self:next(true)
		elseif direction == 'up' then
			self:previous(true)
		end
		return
	end

	local lines = self.pages[self.curPage].lines

	local oldOffset = self.offset
	local oldHorizOffset = self.horizOffset
	if direction == 'down' then
		self.offset = math.min(self.offset + self.step.vertical, math.max(1, #lines - self.region.height))
	elseif direction == 'up' then
		self.offset = math.max(self.offset - self.step.vertical, 1)
	end

--[[
	if direction == 'left' then
		self.horizOffset = math.max(self.horizOffset - self.step.horizontal, 1)
	elseif direction == 'right' then
		self.horizOffset = self.horizOffset + self.step.horizontal
	end
]]--

	if self.offset ~= oldOffset then self:draw() end
	if self.horizOffset ~= oldHorizOffset then self:draw() end
end

function Greenhouse:update()
	self:resize()
	if self.isSpecial then
		self:updateSpecial()
	end

	self:draw()
end


function Greenhouse:special(val)
	self.isSpecial = val
	self:update()
end

function Greenhouse:toggleSpecial()
	self:special(not self.isSpecial)
end

--- This function will be called when the special page
--- is on and needs to be updated.
function Greenhouse:updateSpecial()
end

function Greenhouse:contents()
end

function Greenhouse:toc(toggle)
	if not self.isSpecial then
		self.specialPageIdx = self.curPage
	end
	if toggle then self.isSpecial = not self.isSpecial end
	-- Generate a special page for our table of contents
	local tocText = string.format([[
%s

]], lunacolors.cyan(lunacolors.bold '―― Table of Contents ――'))

	local genericPageCount = 1
	local contents = self:contents()
	if contents then
		for i, c in ipairs(contents) do
			local title = c.title
			if c.active then
				title = lunacolors.invert(title)
			end

			tocText = tocText .. title .. '\n'
		end
	else
		for i, page in ipairs(self.pages) do
			local title = page.title
			if title == 'Page' then
				title = 'Page #' .. genericPageCount
				genericPageCount = genericPageCount + 1
			end
			if i == self.specialPageIdx then
				title = lunacolors.invert(title)
			end

			tocText = tocText .. title .. '\n'
		end
	end
	self.specialPage = Page('TOC', tocText)
	function self:updateSpecial()
		self:toc()
	end
	self:draw()
end

function Greenhouse:resize()
	local size = terminal.size()
	self.region = size
end

function Greenhouse:next(special)
	local oldCurrent = special and self.specialPageIdx or self.curPage
	local pageIdx = math.min(oldCurrent + 1, #self.pages)

	if special then
		self.specialPageIdx = pageIdx
	else
		self.curPage = pageIdx
	end

	if pageIdx ~= oldCurrent then
		self.offset = 1
		self:update()
	end
end

function Greenhouse:previous(special)
	local oldCurrent = special and self.specialPageIdx or self.curPage
	local pageIdx = math.max(self.curPage - 1, 1)

	if special then
		self.specialPageIdx = pageIdx
	else
		self.curPage = pageIdx
	end

	if pageIdx ~= oldCurrent then
		self.offset = 1
		self:update()
	end
end

function Greenhouse:jump(idx)
	if idx ~= self.curPage then
		self.offset = 1
	end
	self.curPage = idx
	self:update()
end

function Greenhouse:keybind(key, callback)
	self.keybinds[key] = callback
end

function Greenhouse:input(char)
end

function Greenhouse:initUi()
	local ansikit = require 'ansikit'
	local bait = require 'bait'
	local commander = require 'commander'
	local hilbish = require 'hilbish'
	local terminal = require 'terminal'
	local Page = require 'nature.greenhouse.page'
	local done = false

	bait.catch('signal.sigint', function()
		ansikit.clear()
		done = true
	end)

	bait.catch('signal.resize', function()
		self:update()
	end)

	ansikit.screenAlt()
	ansikit.clear(true)
	self:draw()

	hilbish.goro(function()
		while not done do
			local c = read()
			self:keybind('Ctrl-D', function()
				done = true
			end)

			if self.keybinds[c] then
				self.keybinds[c](self)
			else
				self:input(c)
			end

	--[[
			if c == 27 then
				local c1 = read()
				if c1 == 91 then
					local c2 = read()
					if c2 == 66 then -- arrow down
						self:scroll 'down'
					elseif c2 == 65 then -- arrow up
						self:scroll 'up'
					end

					if c2 == 49 then
						local c3 = read()
						if c3 == 59 then
							local c4 = read()
							if c4 == 53 then
								local c5 = read()
								if c5 == 67 then
									self:next()
								elseif c5 == 68 then
									self:previous()
								end
							end
						end
					end
				end
				goto continue
			end
			]]--

			::continue::
		end
	end)

	while not done do
		--
	end
	ansikit.showCursor()
	ansikit.screenMain()
end

function read()
	terminal.saveState()
	terminal.setRaw()
	local c = hilbish.editor.readChar()

	terminal.restoreState()
	return c
end

return Greenhouse
