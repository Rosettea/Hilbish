local bait = require 'bait'
local lunacolors = require 'lunacolors'

bait.catch('hilbish.init', function()
	if hilbish.interactive and type(hilbish.opts.greeting) == 'string' then
		if os.date '%d' == '01' and os.date '%m' == '04' then
			print('welcome to a shell, i think??')
		else
			print(lunacolors.format(hilbish.opts.greeting))
		end
	end
end)
