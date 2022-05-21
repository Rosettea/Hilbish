local commander = require 'commander'

commander.register('disown', function(args)
	local id = tonumber(args[1])
	if not id then
		print 'invalid id for job'
		return 1
	end

	local ok = pcall(hilbish.jobs.disown, id)
	if not ok then
		print 'job does not exist'
	end
end)
