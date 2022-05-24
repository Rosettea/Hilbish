local commander = require 'commander'

commander.register('bg', function()
	local job = hilbish.jobs.last()
	if not job then
		print 'bg: no last job'
		return 1
	end

	local err = job.background()
	if err then
		print('bg: ' .. err)
		return 2
	end
end)
