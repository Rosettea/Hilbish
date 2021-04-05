-- Default Hilbish config
ansikit = require 'ansikit'
bait = require 'bait'

function doPrompt(fail)
	prompt(ansikit.format(
		'{blue}%u {cyan}%d ' .. (fail and '{red}' or '{green}') .. '∆{reset} '
	))
end

print(ansikit.format('Welcome {cyan}'.. os.getenv 'USER' ..
'{reset} to {magenta}Hilbish{reset},\n' .. 
'the nice lil shell for {blue}Lua{reset} fanatics!\n'))

doPrompt()

bait.catch('command.exit', function(code)
	doPrompt(code ~= 0)
end)

--hook("tab complete", function ())
