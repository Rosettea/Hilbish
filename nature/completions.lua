function hilbish.completion.handler(line, pos)
	if type(line) ~= 'string' then error '#1 must be a string' end
	if type(pos) ~= 'number' then error '#2 must be a number' end

	-- trim leading whitespace
	local ctx = line:gsub('^%s*(.-)$', '%1')
	if ctx:len() == 0 then return {}, '' end

	ctx = hilbish.aliases.resolve(ctx)
	local fields = string.split(ctx, ' ')
	local query = fields[#fields]

	if #fields == 1 then
		local comps, pfx = hilbish.completion.bins(query, ctx, fields)
		local compGroup = {
			items = comps,
			type = 'grid'
		}

		return {compGroup}, pfx
	else
		local ok, compGroups, pfx = pcall(hilbish.completion.call,
		'command.' .. #fields[1], query, ctx, fields)
		if ok then
			return compGroups, pfx
		end

		local comps, pfx = hilbish.completion.files(query, ctx, fields)
		local compGroup = {
			items = comps,
			type = 'grid'
		}

		return {compGroup}, pfx
	end
end
