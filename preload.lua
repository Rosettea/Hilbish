-- The preload file initializes everything else for our shell
-- Currently it just adds our builtins

local fs = require 'fs'
local commander = require 'commander'
local bait = require 'bait'

commander.register('cd', function (path)
	if #path == 1 then
		local ok, err = pcall(function() fs.cd(path[1]) end)
		if not ok then
			if err == 1 then
				print('directory does not exist')
			end
		end
		bait.throw('cd', path)
		return
	end
	fs.cd(os.getenv 'HOME')
	bait.throw('cd', path)
end)
