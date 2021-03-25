-- Default Hilbish config
ansikit = require 'ansikit'

prompt(ansikit.text(
	'{blue}%u {cyan}%d {green}∆{reset} '
))

--hook("tab complete", function ())
