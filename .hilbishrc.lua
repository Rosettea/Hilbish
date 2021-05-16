-- Default Hilbish config
lunacolors = require 'lunacolors'
bait = require 'bait'

function doPrompt(fail)
	prompt(lunacolors.format(
		'{blue}%u {cyan}%d ' .. (fail and '{red}' or '{green}') .. '∆ '
	))
end

print(lunacolors.format('Welcome to {magenta}Hilbish{reset}, {cyan}' .. hilbish.user
.. '{reset}.\n' .. 'The nice lil shell for {blue}Lua{reset} fanatics!\n'))

doPrompt()

bait.catch('command.exit', function(code)
	doPrompt(code ~= 0)
end)

