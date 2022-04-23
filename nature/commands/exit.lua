local bait = require 'bait'
local commander = require 'commander'

commander.register('exit', function()
	bait.throw('hilbish.exit')
	os.exit(0)
end)
