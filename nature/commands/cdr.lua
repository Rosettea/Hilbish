local commander = require 'commander'
local fs = require 'fs'
local lunacolors = require 'lunacolors'
local dirs = require 'nature.dirs'

commander.register('cdr', function(args)
	if not args[1] then
		print(lunacolors.format [[
cdr: change directory to one which has been recently visied

usage: cdr <index>

to get a list of recent directories, use {green}{underline}cdr list{reset}]])
		return
	end

	if args[1] == 'list' then
		local recentDirs = dirs.recentDirs
		if #recentDirs == 0 then
			print 'No directories have been visited.'
			return 1
		end
		print(table.concat(recentDirs, '\n'))
		return
	end

	local index = tonumber(args[1])
	if not index then
		print(string.format('Received %s as index, which isn\'t a number.', index))
		return 1
	end

	if not dirs.recent(index) then
		print(string.format('No recent directory found at index %s.', index))
		return 1
	end

	fs.cd(dirs.recent(index))
end)
