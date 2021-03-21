-- Default Hilbish config
ansikit = require 'ansikit'

prompt(ansikit.text('λ {bold}{cyan}'..os.getenv('USER')..' >{magenta}>{cyan}>{reset} '))

--hook("tab complete", function ())
