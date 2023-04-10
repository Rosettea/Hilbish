local bait = require 'bait'
local lunacolors = require 'lunacolors'

bait.catch('hilbish.init', function()
	if hilbish.interactive and type(hilbish.opts.greeting) == 'string' then
		print(lunacolors.format(hilbish.opts.greeting))
	end
end)
