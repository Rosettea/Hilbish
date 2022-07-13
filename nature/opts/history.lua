local bait = require 'bait'

bait.catch('command.exit', function(_, cmd, priv)
	if not priv and hilbish.opts.history then hilbish.history.add(cmd) end
end)
