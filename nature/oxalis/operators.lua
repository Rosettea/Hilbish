local fs = require 'fs'
local M = {}

function M.pipe(cmd, cmd2)
	local pr, pw = fs.pipe()
	hilbish.run(cmd, {
		out = pw,
		err = pw,
	})
	pw:close()
	hilbish.run(cmd2, {
		input = pr
	})
	return {command = cmd2, input = pr, out }
end

return M
