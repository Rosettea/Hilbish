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
