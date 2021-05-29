-- The preload file initializes everything else for our shell

local fs = require 'fs'
local commander = require 'commander'
local bait = require 'bait'
local old_dir = hilbish.cwd()

local shlvl = tonumber(os.getenv 'SHLVL')
if shlvl ~= nil then os.setenv('SHLVL', shlvl + 1) else os.setenv('SHLVL', 1) end

-- Builtins
commander.register('cd', function (args)
	bait.throw('cd', args)
	if #args > 0 then
		local path = ''
		for i = 1, #args do
			path = path .. tostring(args[i]) .. ' '
		end
		path = path:gsub('$%$','\0'):gsub('${([%w_]+)}', os.getenv)
		:gsub('$([%w_]+)', os.getenv):gsub('%z','$'):gsub("%s+", "")

        if path == '-' then
            path = old_dir
            print(path)
        end
        old_dir = hilbish.cwd()

		local ok, err = pcall(function() fs.cd(path) end)
		if not ok then
			if err == 1 then
				print('directory does not exist')
			end
			return err
		end
		return
	end
	fs.cd(os.getenv 'HOME')
	bait.throw('command.exit', 0)
end)

commander.register('exit', function()
	os.exit(0)
end)

do
	local virt_G = { }
	
	setmetatable(_G, {
		__index = function (self, key)
			local got_virt = virt_G[key]
			if got_virt ~= nil then
				return got_virt
			end
			
			virt_G[key] = os.getenv(key)
			return virt_G[key]
		end,
			
		__newindex = function (self, key, value)
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

-- Function additions to Lua standard library
function string.split(str, delimiter)
	local result = {}
	local from = 1
	
	local delim_from, delim_to = string.find(str, delimiter, from)
	
	while delim_from do
		table.insert(result, string.sub(str, from, delim_from - 1))
		from = delim_to + 1
		delim_from, delim_to = string.find(str, delimiter, from)
	end
	
	table.insert(result, string.sub(str, from))

	return result
end

-- Hook handles
bait.catch('command.not-found', function(cmd)
	print(string.format('hilbish: %s not found', cmd))
end)

