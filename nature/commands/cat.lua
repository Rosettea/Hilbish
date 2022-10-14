local commander = require 'commander'
local fs = require 'fs'

commander.register('cat', function(args)
	local exit = 0

	if #args == 0 then
		print [[
usage: cat [file]...]]
	end

	for _, fName in ipairs(args) do
		local f = io.open(fName)
		if f == nil then
			exit = 1
			print(string.format('cat: %s: no such file or directory', fName))
			goto continue
		end

		io.write(f:read '*a')
		::continue::
	end
	io.flush()
	return exit
end)
