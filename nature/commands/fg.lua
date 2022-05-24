local commander = require 'commander'

commander.register('fg', function()
	local job = hilbish.jobs.last()
	if not job then
		print 'fg: no last job'
		return 1
	end

	local err = job.foreground() -- waits for job; blocks
	if err then
		print('fg: ' .. err)
		return 2
	end
end)
