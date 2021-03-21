-- The preload file initializes everything else for our shell
-- Currently it just adds our builtins

local fs = require 'fs'
local commander = require 'commander'

commander.register('cd', function (path)
	if #path == 1 then
		fs.cd(path[1])
	end
end)
