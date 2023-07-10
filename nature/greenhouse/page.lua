local Object = require 'nature.object'

local Page = Object:extend()

function Page:new(title, text)
	self:setText(text)
	self.title = title or 'Page'
end

function Page:setText(text)
	self.lines = string.split(text, '\n')
end

function Page:setTitle(title)
	self.title = title
end

return Page
