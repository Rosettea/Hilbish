local ansikit = require 'ansikit'
local bait = require 'bait'
local commander = require 'commander'
local hilbish = require 'hilbish'
local terminal = require 'terminal'
local Greenhouse = require 'nature.greenhouse'
local Page = require 'nature.greenhouse.page'

commander.register('greenhouse', function(args, sinks)
	local gh = Greenhouse(sinks.out)
	function gh:resize()
		local size = terminal.size()
		self.region = {
			width = size.width,
			height = size.height - 2
		}
	end
	local oldDraw = gh.draw
	function gh:draw()
		oldDraw(self)
		local workingPage = self.pages[self.curPage]
		local offset = self.offset
		if self.isToc then
			offset = self.tocOffset
			workingPage = self.tocPage
		end

		self.sink:write(ansikit.getCSI((self.region.height + 2) - self.start.. ';1', 'H'))
		if not self.isToc then
			self.sink:write(string.format('\27[0mPage %d', self.curPage))
			if workingPage.title ~= '' then
				self.sink:writeln(' — ' .. workingPage.title)
			else
				self.sink:writeln('')
			end
		end
		self.sink:write(buffer)
	end
	gh:resize()

	gh:keybind('Enter', function(self)
		if self.isToc then
			self:jump(self.tocPageIdx)
			self:toc(true)
		end
	end)

	if sinks['in'].pipe then
		local page = Page('', sinks['in']:readAll())
		gh:addPage(page)
	end

	for _, name in ipairs(args) do
		local f <close> = io.open(name, 'r')
		if not f then
			sinks.err:writeln(string.format('could not open file %s', name))
		end

		local page = Page(name, f:read '*a')
		gh:addPage(page)
	end

	gh:initUi()
end)
