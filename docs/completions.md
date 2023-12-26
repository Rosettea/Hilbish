---
title: Completions
description: Tab completion for commands.
layout: doc
menu: 
  docs:
    parent: "Features"
---

Completions for commands can be created with the [`hilbish.complete`](../api/hilbish#complete)
function. See the link for how to use it.

To create completions for a command is simple.
The callback will be passed 3 parameters:
- `query` (string): The text that the user is currently trying to complete.
This should be used to match entries.
- `ctx` (string): Contains the entire line. Use this if
more text is needed to be parsed for context.
- `fields` (string): The `ctx` split up by spaces.

In most cases, the completer just uses `fields` to check the amount
and `query` on what to match entries on.

In order to return your results, it has to go within a "completion group."
Then you return a table of completion groups and a prefix. The prefix will
usually just be the `query`.

Hilbish allows one to mix completion menus of different types, so
a grid menu and a list menu can be used and complete and display at the same time.
A completion group is a table with these keys:
- `type` (string): type of completion menu, either `grid` or `list`.
- `items` (table): a list of items. 

The requirements of the `items` table is different based on the
`type`. If it is a `grid`, it can simply be a table of strings.

Otherwise if it is a `list` then each entry can
either be a string or a table.
Example:
```lua
local cg = {
	items = {
		'list item 1',
		['--command-flag-here'] = {'this does a thing', '--the-flag-alias'}
	},
	type = 'list'
}
local cg2 = {
	items = {'just', 'a bunch', 'of items', 'here', 'hehe'},
	type = 'grid'
}

return {cg, cg2}, prefix
```

Which looks like this:  
{{< video src="https://safe.saya.moe/t4CiLK6dgPbD.mp4" >}}

# Completion Handler
Like most parts of Hilbish, it's made to be extensible and
customizable. The default handler for completions in general can
be overwritten to provide more advanced completions if needed.
This usually doesn't need to be done though, unless you know
what you're doing.

The default completion handler provides 3 things:
binaries (with a plain name requested to complete, those in
$PATH), files, or command completions. It will try to run a handler
for the  command or fallback to file completions.

To overwrite it, just assign a function to `hilbish.completion.handler` like so:
```lua
-- line is the entire line as a string
-- pos is the position of the cursor.
function hilbish.completion.handler(line, pos)
	-- do things
end
```
