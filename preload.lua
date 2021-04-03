-- The preload file initializes everything else for our shell
-- Currently it just adds our builtins

local fs = require 'fs'
local commander = require 'commander'
local bait = require 'bait'

commander.register('cd', function (args)
	bait.throw('cd', args)
	if #args > 0 then
		local path = ''
		for i = 1, #args do
			path = path .. tostring(args[i]) .. ' '
		end

		local ok, err = pcall(function() fs.cd(path) end)
		if not ok then
			if err == 1 then
				print('directory does not exist')
			end
			bait.throw('command.fail', nil)
		else bait.throw('command.success', nil) end
		return
	end
	fs.cd(os.getenv 'HOME')
	bait.throw('command.success', nil)
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
