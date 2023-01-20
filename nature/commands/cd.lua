local bait = require 'bait'
local commander = require 'commander'
local fs = require 'fs'
local dirs = require 'nature.dirs'

dirs.old = hilbish.cwd()
commander.register('cd', function (args, sinks)
	if #args > 1 then
		sinks.out:writeln("cd: too many arguments")
		return 1
	elseif #args > 0 then
		local path = args[1]:gsub('$%$','\0'):gsub('${([%w_]+)}', os.getenv)
		:gsub('$([%w_]+)', os.getenv):gsub('%z','$'):gsub('^%s*(.-)%s*$', '%1')

        if path == '-' then
			path = dirs.old
            sinks.out:writeln(path)
        end
        dirs.setOld(hilbish.cwd())
		dirs.push(path)

		local ok, err = pcall(function() fs.cd(path) end)
		if not ok then
			sinks.out:writeln(err)
			return 1
		end
		bait.throw('cd', path)

		return
	end
	fs.cd(hilbish.home)
	bait.throw('cd', hilbish.home)

	dirs.push(hilbish.home)
end)
