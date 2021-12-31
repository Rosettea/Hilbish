-- Default Hilbish config
lunacolors = require 'lunacolors'
bait = require 'bait'

function doPrompt(fail)
	prompt(lunacolors.format(
		'{blue}%u {cyan}%d ' .. (fail and '{red}' or '{green}') .. '∆ '
	))
end

print(lunacolors.format(hilbish.greeting))

doPrompt()

bait.catch('command.exit', function(code)
	doPrompt(code ~= 0)
end)

