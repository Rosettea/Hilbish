defmodule Parse do
	def md_to_html(a) do
		import MDEx
		MDEx.to_html!(a, extension: [
			strikethrough: true,
			tagfilter: true,
			table: true,
			autolink: true,
			tasklist: true,
			footnotes: true,
			shortcodes: true,
		],
		parse: [
			smart: true,
			relaxed_tasklist_matching: true,
			relaxed_autolinks: true
		],
		render: [
			github_pre_lang: true,
			unsafe_: true,
		],
		features: [
			sanitize: true
		])
	end
end
