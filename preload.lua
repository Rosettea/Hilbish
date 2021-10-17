-- The preload file initializes everything else for our shell

local fs = require 'fs'
local commander = require 'commander'
local bait = require 'bait'
require 'succulent' -- Function additions
local oldDir = hilbish.cwd()

local shlvl = tonumber(os.getenv 'SHLVL')
if shlvl ~= nil then os.setenv('SHLVL', shlvl + 1) else os.setenv('SHLVL', 1) end

-- Builtins
local recentDirs = {}
commander.register('cd', function (args)
	if #args > 0 then
		local path = table.concat(args, ' '):gsub('$%$','\0'):gsub('${([%w_]+)}', os.getenv)
		:gsub('$([%w_]+)', os.getenv):gsub('%z','$'):gsub('^%s*(.-)%s*$', '%1')

        if path == '-' then
            path = oldDir
            print(path)
        end
        oldDir = hilbish.cwd()

		local ok, err = pcall(function() fs.cd(path) end)
		if not ok then
			print(err:sub(17))
			return 1
		end
		bait.throw('cd', path)

		-- add to table of recent dirs
		table.insert(recentDirs, 1, path)
		recentDirs[11] = nil

		return
	end
	fs.cd(hilbish.home)
	bait.throw('cd', hilbish.home)

	return
end)

commander.register('exit', function()
	os.exit(0)
end)

commander.register('doc', function(args)
	local moddocPath = hilbish.dataDir .. '/docs/'
	local globalDesc = [[
These are the global Hilbish functions that are always available and not part of a module.]]
	if #args > 0 then
		local mod = table.concat(args, ' '):gsub('^%s*(.-)%s*$', '%1')

		local f = io.open(moddocPath .. mod .. '.txt', 'rb')
		if not f then 
			print('Could not find docs for module named ' .. mod .. '.')
			return 1
		end

		local desc = (mod == 'global' and globalDesc or getmetatable(require(mod)).__doc)
		local funcdocs = f:read '*a'
		local backtickOccurence = 0
		print(desc .. '\n\n' .. lunacolors.format(funcdocs:sub(1, #funcdocs - 1):gsub('`', function()
			backtickOccurence = backtickOccurence + 1
			if backtickOccurence % 2 == 0 then
				return '{reset}'
			else
				return '{underline}{green}'
			end
		end)))
		f:close()

		return
	end
	local modules = table.map(fs.readdir(moddocPath), function(f)
		return lunacolors.underline(lunacolors.blue(f:sub(1, -5)))
	end)

	io.write [[
Welcome to Hilbish's doc tool! Here you can find documentation for builtin
functions and other things.

Usage: doc <module>

Available modules: ]]

	print(table.concat(modules, ', '))

	return
end)

do
	local virt_G = { }

	setmetatable(_G, {
		__index = function (_, key)
			local got_virt = virt_G[key]
			if got_virt ~= nil then
				return got_virt
			end

			virt_G[key] = os.getenv(key)
			return virt_G[key]
		end,

		__newindex = function (_, key, value)
			if type(value) == 'string' then
				os.setenv(key, value)
				virt_G[key] = value
			else
				if type(virt_G[key]) == 'string' then
					os.setenv(key, '')
				end
				virt_G[key] = value
			end
		end,
	})

	bait.catch('command.exit', function ()
		for key, value in pairs(virt_G) do
			if type(value) == 'string' then
				virt_G[key] = os.getenv(key)
			end
		end
	end)
end

commander.register('cdr', function(args)
	if not args[1] then
		print(lunacolors.format [[
cdr: change directory to one which has been recently visied

usage: cdr <index>

to get a list of recent directories, use {green}{underline}cdr list{reset}]])
		return
	end

	if args[1] == 'list' then
		if #recentDirs == 0 then
			print 'No directories have been visited.'
			return 1
		end
		print(table.concat(recentDirs, '\n'))
		return
	end

	local index = tonumber(args[1])
	if not index then
		print(string.format('received %s as index, which isn\'t a number', index))
		return 1
	end

	if not recentDirs[index] then
		print(string.format('no recent directory found at index %s', index))
		return 1
	end

	fs.cd(recentDirs[index])
	return
end)

-- Hook handles
bait.catch('command.not-found', function(cmd)
	print(string.format('hilbish: %s not found', cmd))
end)

