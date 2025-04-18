-- @module hilbish.processors

hilbish.processors = {
	list = {},
	sorted = {}
}

function hilbish.processors.add(processor)
	if not processor.func then
		error 'processor is missing function'
	end

	table.insert(hilbish.processors.list, processor)
	table.sort(hilbish.processors.list, function(a, b) return a.priority < b.priority end)
end

--- Run all command processors, in order by priority.
--- It returns the processed command (which may be the same as the passed command)
--- and a boolean which states whether to proceed with command execution.
function hilbish.processors.execute(command)
	local continue = true
	for _, processor in ipairs(hilbish.processors.list) do
		local processed = processor.func(command)
		if processed.command then command = processed.command end
		if not processed.continue then
			continue = false
			break
		end
	end

	return command, continue
end
