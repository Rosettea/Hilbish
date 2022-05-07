local commander = require 'commander'
local fs = require 'fs'
local lunacolors = require 'lunacolors'

commander.register('doc', function(args)
	local moddocPath = hilbish.dataDir .. '/docs/'
	local modDocFormat = [[
%s
%s
# Functions
]]

	if #args > 0 then
		local mod = args[1]

		local f = io.open(moddocPath .. mod .. '.txt', 'rb')
		local funcdocs = nil
		if not f then
			-- assume subdir
			-- dataDir/docs/<mod>/<mod>.txt
			moddocPath = moddocPath .. mod .. '/'
			local subdocName = args[2]
			if not subdocName then
				subdocName = 'index'
			end
			f = io.open(moddocPath .. subdocName .. '.txt', 'rb')
			if not f then
				print('No documentation found for ' .. mod .. '.')
				return
			end
			funcdocs = f:read '*a'
			local moddocs = table.filter(fs.readdir(moddocPath), function(f) return f ~= 'index.txt' end)
			local subdocs = table.map(moddocs, function(fname)
				return lunacolors.underline(lunacolors.blue(string.gsub(fname, '.txt', '')))
			end)
			if subdocName == 'index' then
				funcdocs = funcdocs .. '\nSubdocs: ' .. table.concat(subdocs, ', ')
			end
		end

		if not funcdocs then
			funcdocs = f:read '*a'
		end
		local desc = ''
		local ok = pcall(require, mod)
		local backtickOccurence = 0
		local formattedFuncs = lunacolors.format(funcdocs:sub(1, #funcdocs - 1):gsub('`', function()
			backtickOccurence = backtickOccurence + 1
			if backtickOccurence % 2 == 0 then
				return '{reset}'
			else
				return '{underline}{green}'
			end
		end))

		if ok then
			local props = {}
			local propstr = ''
			local modDesc = ''
			local modmt = getmetatable(require(mod))
			modDesc = modmt.__doc
			if modmt.__docProp then
				-- not all modules have docs for properties
				props = table.map(modmt.__docProp, function(v, k)
					return lunacolors.underline(lunacolors.blue(k)) .. ' > ' .. v
				end)
			end
			if #props > 0 then
				propstr = '\n# Properties\n' .. table.concat(props, '\n') .. '\n'
			end
			desc = string.format(modDocFormat, modDesc, propstr)
		end
		print(desc .. formattedFuncs)
		f:close()

		return
	end
	local modules = table.map(fs.readdir(moddocPath), function(f)
		return lunacolors.underline(lunacolors.blue(string.gsub(f, '.txt', '')))
	end)

	io.write [[
Welcome to Hilbish's doc tool! Here you can find documentation for builtin
functions and other things.

Usage: doc <section> [subdoc]
A section is a module or a literal section and a subdoc is a subsection for it.

Available sections: ]]
	io.flush()

	print(table.concat(modules, ', '))
end)
