-- Default Hilbish config
ansikit = require 'ansikit'
bait = require 'bait'

function doPrompt(fail)
	prompt(ansikit.text(
		'{blue}%u {cyan}%d ' .. (fail and '{red}' or '{green}') .. '∆{reset} '
	))
end

print(ansikit.text('Welcome {cyan}'.. os.getenv 'USER' ..
'{reset} to {magenta}Hilbish{reset},\n' .. 
'the nice lil shell for {blue}Lua{reset} fanatics!\n'))

doPrompt()

bait.catch('command.fail', function()
	doPrompt(true)
end)

bait.catch('command.success', function()
	doPrompt()
end)

--hook("tab complete", function ())
