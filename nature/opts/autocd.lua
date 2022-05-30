local fs = require 'fs'

function cdHandle(inp)
	local res = hilbish.runner.lua(inp)

	if not res.err then
		return res
	end

	res = hilbish.runner.sh(inp)

	if res.exit ~= 0 and hilbish.opts.autocd then
		local ok, stat = pcall(fs.stat, res.input)
		if ok and stat.isDir then
			-- discard here to not append the cd, which will be in history
			local _, exitCode, err = hilbish.runner.sh('cd ' .. res.input)
			res.exitCode = exitCode
			res.err = err
		end
	end

	return res
end

hilbish.runner.setMode(cdHandle)
