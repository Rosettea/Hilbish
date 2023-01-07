local commander = require 'commander'
local fs = require 'fs'
local lunacolors = require 'lunacolors'

commander.register('doc', function(args)
	local moddocPath = hilbish.dataDir .. '/docs/'
	local apidocHeader = [[
# %s
{grayBg}  {white}{italic}%s  {reset}

]]

	local modules = table.map(fs.readdir(moddocPath), function(f)
		return lunacolors.underline(lunacolors.blue(string.gsub(f, '.md', '')))
	end)
	local doc = [[
Welcome to Hilbish's documentation viewer! Here you can find
documentation for builtin functions and other things related
to Hilbish.

Usage: doc <section> [subdoc]
Available sections: ]] .. table.concat(modules, ', ')
	if #args > 0 then
		local mod = args[1]

		local f = io.open(moddocPath .. mod .. '.md', 'rb')
		local funcdocs = nil
		local subdocName = args[2]
		if not f then
			-- assume subdir
			-- dataDir/docs/<mod>/<mod>.md
			moddocPath = moddocPath .. mod .. '/'
			if not subdocName then
				subdocName = '_index'
			end
			f = io.open(moddocPath .. subdocName .. '.md', 'rb')
			if not f then
				moddocPath = moddocPath .. subdocName .. '/'
				subdocName = args[3] or '_index'
				f = io.open(moddocPath .. subdocName .. '.md', 'rb')
			end
			if not f then
				print('No documentation found for ' .. args[#args] .. '.')
				return
			end
		end
		funcdocs = f:read '*a':gsub('-([%d]+)', '%1')
		local moddocs = table.filter(fs.readdir(moddocPath), function(f) return f ~= '_index.md' end)
		local subdocs = table.map(moddocs, function(fname)
			return lunacolors.underline(lunacolors.blue(string.gsub(fname, '.md', '')))
		end)
		if subdocName == '_index' then
			funcdocs = funcdocs .. '\nSubdocs: ' .. table.concat(subdocs, ', ')
		end

		local valsStr = funcdocs:match '%-%-%-\n([^%-%-%-]+)\n'
		local vals = {}
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
		if mod == 'api' then
			funcdocs = string.format(apidocHeader, vals.title, vals.description or 'no description.') .. funcdocs
		end
		doc = funcdocs:sub(1, #funcdocs - 1)
		f:close()
	end

	local backtickOccurence = 0
	print(lunacolors.format(doc:gsub('`', function()
		backtickOccurence = backtickOccurence + 1
		if backtickOccurence % 2 == 0 then
			return '{reset}'
		else
			return '{underline}{green}'
		end
	end):gsub('#+.-\n', function(t)
		local signature = t:gsub('<.->(.-)</.->', '{underline}%1'):gsub('\\', '<')
		return '{bold}{yellow}' .. signature .. '{reset}'
	end)))
end)
