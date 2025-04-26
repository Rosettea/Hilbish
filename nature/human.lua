local fs = require 'fs'

local M = {}

-- Find where a manpage is.
function M.where(name, sections)
	local manpath = os.getenv 'MANPATH'
	if not manpath then
		manpath = '/usr/local/share/man:/usr/share/man'
	end

	local paths = string.split(manpath, ':')
	for _, path in ipairs(paths) do
		-- man directory structure:
		-- <manpath>/man[sectionNumber]/manpage.[sectionNumber].gz
		-- example: <manpath>/man1/Xorg.1.gz
		local manSubPaths = fs.glob(fs.join(path, string.format('man%s/', sections and '[' .. table.concat(sections, '') .. ']' or '*')))
		for _, subPath in ipairs(manSubPaths) do
			local currentSection = subPath:match '/man([%w%d]+)$'
			local assumedPath = fs.join(subPath, string.format('%s%s', name, '.*' .. currentSection .. '.gz'))
			local globbedPages = fs.glob(assumedPath)
			if globbedPages[1] then
				return globbedPages[1]
			end
		end
	end
end

function M.parse(path)
	assert(fs.stat(path), 'file does not exist')

	local _, contents = hilbish.run(string.format('gzip -d %s -c', path), false)
	local sections = {}
	local sectionPattern = '\n%.SH%s+([^\n]+)\n(.-)\n%.SH'
	for sectionName, sectionContent in string.gmatch(contents, sectionPattern) do
		sections[string.lower(sectionName)] = sectionContent
	end

	return sections
end

return M
