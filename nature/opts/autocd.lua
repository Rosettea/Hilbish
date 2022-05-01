local fs = require 'fs'

function cdHandle(inp)
	local input, exit, err = hilbish.runner.lua(inp)

	if not err then
		return input, exit, err
	end

	input, exit, err = hilbish.runner.sh(inp)

	if exit ~= 0 and hilbish.opts.autocd then
		local ok, stat = pcall(fs.stat, input)
		if ok and stat.isDir then
			-- discard here to not append the cd, which will be in history
			_, exit, err = hilbish.runner.sh('cd ' .. input)
		end
	end

	return input, exit, err
end

hilbish.runner.setMode(cdHandle)
