local commander = require 'commander'

commander.register('exec', function(args)
	hilbish.exec(args[1])
end)
