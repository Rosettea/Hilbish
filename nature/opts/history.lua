local bait = require 'bait'

bait.catch('command.exit', function(_, cmd, priv)
	if not cmd then return end
	if not priv and hilbish.opts.history then hilbish.history.add(cmd) end
end)
