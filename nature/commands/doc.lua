local ansikit = require 'ansikit'
local commander = require 'commander'
local fs = require 'fs'
local lunacolors = require 'lunacolors'
local Greenhouse = require 'nature.greenhouse'
local Page = require 'nature.greenhouse.page'

commander.register('doc', function(args, sinks)
	local moddocPath = hilbish.dataDir .. '/docs/'
	local stat = pcall(fs.stat, '.git/refs/heads/extended-job-api')
	if stat then
		-- hilbish git
		moddocPath = './docs/'
	end

	local modules = table.map(fs.readdir(moddocPath), function(f)
		return lunacolors.underline(lunacolors.blue(string.gsub(f, '.md', '')))
	end)
	local doc = [[
Welcome to Hilbish's documentation viewer! Here you can find
documentation for builtin functions and other things related
to Hilbish.

Usage: doc <section> [subdoc]
Available sections: ]] .. table.concat(modules, ', ')
	local vals = {}

	if #args > 0 then
		local mod = args[1]

		local f = io.open(moddocPath .. mod .. '.md', 'rb')
		local funcdocs = nil
		local subdocName = args[2]
		if not f then
			moddocPath = moddocPath .. mod .. '/'
			if not subdocName then
				subdocName = '_index'
			end
			f = io.open(moddocPath .. subdocName .. '.md', 'rb')
			local oldmoddocPath = moddocPath
			if not f then
				moddocPath = moddocPath .. subdocName:match '%w+' .. '/'
				f = io.open(moddocPath .. subdocName .. '.md', 'rb')
			end
			if not f then
				moddocPath = oldmoddocPath .. subdocName .. '/'
				subdocName = args[3] or '_index'
				f = io.open(moddocPath .. subdocName .. '.md', 'rb')
			end
			if not f then
				sinks.out:writeln('No documentation found for ' .. mod .. '.')
				return 1
			end
		end
		funcdocs = f:read '*a':gsub('-([%d]+)', '%1')
		local moddocs = table.filter(fs.readdir(moddocPath), function(f) return f ~= '_index.md' and f ~= 'index.md' end)
		local subdocs = table.map(moddocs, function(fname)
			return lunacolors.underline(lunacolors.blue(string.gsub(fname, '.md', '')))
		end)
		if #moddocs ~= 0 then
			funcdocs = funcdocs .. '\nSubdocs: ' .. table.concat(subdocs, ', ') .. '\n\n'
		end

		local valsStr = funcdocs:match '%-%-%-\n([^%-%-%-]+)\n'
		if valsStr then
			local _, endpos = funcdocs:find('---\n' .. valsStr .. '\n---\n\n', 1, true)
			funcdocs = funcdocs:sub(endpos + 1, #funcdocs)

			-- parse vals
			local lines = string.split(valsStr, '\n')
			for _, line in ipairs(lines) do
				local key = line:match '(%w+): '
				local val = line:match '^%w+: (.-)$'

				if key then
					vals[key] = val
				end
			end
		end
		doc = funcdocs:sub(1, #funcdocs - 1)
		f:close()
	end

	local gh = Greenhouse(sinks.out)
	function gh:resize()
		local size = terminal.size()
		self.region = {
			width = size.width,
			height = size.height - 3
		}
	end
	gh:resize()

	function gh:render()
		local workingPage = self.pages[self.curPage]
		local offset = self.offset
		if self.isSpecial then
			offset = self.specialOffset
			workingPage = self.specialPage
		end

		self.sink:write(ansikit.getCSI(self.region.height + 2 .. ';1', 'H'))
		if not self.isSpecial then
			if args[1] == 'api' then
				self.sink:writeln(lunacolors.reset(string.format('%s', vals.title)))
				self.sink:write(lunacolors.format(string.format('{grayBg} ↳ {white}{italic}%s  {reset}', vals.description or 'No description.')))
			else
				self.sink:write(lunacolors.reset(string.format('Viewing doc page %s', moddocPath)))
			end
		end
	end
	local backtickOccurence = 0
	local page = Page(nil, lunacolors.format(doc:gsub('`', function()
		backtickOccurence = backtickOccurence + 1
		if backtickOccurence % 2 == 0 then
			return '{reset}'
		else
			return '{underline}{green}'
		end
	end):gsub('\n#+.-\n', function(t)
		local signature = t:gsub('<.->(.-)</.->', '{underline}%1'):gsub('\\', '<')
		return '{bold}{yellow}' .. signature .. '{reset}'
	end)))
	gh:addPage(page)
	ansikit.hideCursor()
	gh:initUi()
end)
