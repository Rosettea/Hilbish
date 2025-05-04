import lustre/attribute
import lustre/element
import lustre/element/html

pub fn page() -> element.Element(a) {
	html.div([attribute.class("flex flex-col gap-4")], [
		html.div([attribute.class("border-b border-b-zinc-300 gap-3 -mx-4 p-2 h-screen bg-radial-[at_100%_100%] from-pink-500 to-stone-50 dark:to-neutral-950 to-35% flex flex-col items-center justify-center")], [
			html.div([attribute.class("gap-1 flex flex-col items-center")], [
				html.span([attribute.class("flex flex-row items-center justify-center")], [
					html.img([
						attribute.src("./hilbish-flower.png"),
						attribute.class("h-20")
					]),
					html.p([attribute.class("text-4xl font-bold")], [element.text("Hilbish")]),
				]),
				html.p([attribute.class("text-4xl font-light")], [element.text("Something Unique.")]),
			]),
			html.p([attribute.class("text-center")], [
				element.text("Hilbish is the new Moon-powered interactive shell for Lua fans!"),
				html.br([]),
				element.text("Extensible, scriptable, configurable: All in Lua.")
			]),
			html.div([attribute.class("flex flex-row gap-2 mt-2")], [
				button("Install", "bg-pink-500/30", ""),
				button("GitHub", "bg-stone-500/30", ""),
			]),
			
			html.p([attribute.class("absolute bottom-4")], [element.text("Scroll for more")])
		]),
		html.div([attribute.class("text-center")], [
			html.span([attribute.class("rounded-md backdrop-blur-md bg-pink-500/20 p-2 text-xs font-light")], [element.text("Feature Overview")]),
			html.div([attribute.class("flex flex-col justify-center items-center gap-6 mt-4")], [
				feature_section(
					"The Moon-powered shell",
					"https://upload.wikimedia.org/wikipedia/commons/thumb/c/cf/Lua-Logo.svg/2048px-Lua-Logo.svg.png",
					"Hilbish makes use of the Lua programming language for interactive and config scripting.
					If you write Lua on a regular basis, Hilbish will be the perfect resident in your terminal.
					
					You can still use shell script, but Lua takes the spotlight (or the moonlight..)",
					"start"
				),
				feature_section(
					"Modern, Helpful Interactive Features",
					"https://safe.saya.moe/osR0bplExBC0.png",
					"Graphical TUI history, sensible tab completions, elegantly refreshing prompts, name it all and Hilbish either has it or it's 1 suggestion or 1 script away from being possible.
					Everything present in Hilbish is meant to enhance your interactive shell experience.",
					"end"
				),
				feature_section(
					"Sensible, Friendly Defaults",
					"https://safe.saya.moe/7ze8NQVPD9vO.png",
					"Hilbish's default config makes a simple but presentable showcase of its Lua API and a few of its features.",
					"start"
				),
				feature_section(
					"Truly Make It Yours",
					"",
					"Many things about Hilbish are designed to be changed and swapped out.
					If you want to make use of a Lua derivative in your interactive use (like Fennel) instead of
					Lua, that can be done!",
					"end"
				)
			]),
		]),
		html.div([attribute.class("text-center bg-neutral-100 dark:bg-neutral-900 -mx-4 p-4")], [
			html.span([attribute.class("rounded-md backdrop-blur-md bg-blue-500/20 p-2 text-xs font-light")], [element.text("Download It Now!")]),
			html.div([attribute.class("flex flex-col items-center mt-4 gap-2")], [
				html.p([attribute.class("md:w-3/6")], [element.text("To find out all that Hilbish can do, you should just try it out! It's officially available on Linux, MacOS, Windows, and probably builds on anything Go is available on!")]),
				html.h2([attribute.class("text-3xl font-semibold")], [element.text("Featured Downloads")]),
				html.p([], [element.text("These are binary releases of Hilbish from GitHub. Install instructions are located at an INSTALL.md file.")]),
				html.div([attribute.class("mt-6 flex flex-row flex-wrap items-center justify-center gap-8")], [
					html.div([attribute.class("flex flex-col gap-2")], [
						html.img([
							attribute.src("https://upload.wikimedia.org/wikipedia/commons/thumb/3/35/Tux.svg/1200px-Tux.svg.png"),
							attribute.class("h-36")
						]),
						button("Linux (64-bit)", "bg-stone-500/30", "")
					]),
					html.div([attribute.class("flex flex-col gap-2")], [
						html.img([
							attribute.src("https://upload.wikimedia.org/wikipedia/commons/thumb/0/0a/Unofficial_Windows_logo_variant_-_2002%E2%80%932012_%28Multicolored%29.svg/2321px-Unofficial_Windows_logo_variant_-_2002%E2%80%932012_%28Multicolored%29.svg.png"),
							attribute.class("h-36")
						]),
						button("Windows (64-bit)", "bg-stone-500/30", "")
					]),
					html.div([attribute.class("flex flex-col gap-2 justify-center items-center")], [
						html.img([
							attribute.src("https://images.seeklogo.com/logo-png/38/2/apple-mac-os-logo-png_seeklogo-381401.png"),
							attribute.class("h-36")
						]),
						button("MacOS (64-bit)", "bg-stone-500/30", "")
					]),
					html.div([attribute.class("flex flex-col gap-2 justify-center items-center")], [
						html.img([
							attribute.src("https://images.seeklogo.com/logo-png/38/2/apple-mac-os-logo-png_seeklogo-381401.png"),
							attribute.class("h-36")
						]),
						button("MacOS (ARM)", "bg-stone-500/30", "")
					])
				])
			])
		]),
	])
}

fn feature_section(title: String, image: String, text: String, align: String) -> element.Element(a) {
	let reverse = case align {
		"end" -> "flex-row-reverse"
		_ -> ""
	}
	// for tailwind to generate these styles
	// xl:items-end xl:items-start
	html.div([attribute.class("flex flex-col gap-2 md:w-3/6 text-start xl:items-" <> align)], [
		html.h1([attribute.class("text-4xl font-semibold")], [element.text(title)]),
		html.div([attribute.class("flex flex-row flex-wrap xl:flex-nowrap justify-center items-center gap-4 " <> reverse)], [
			case image {
				"" -> element.none()
				_ -> html.img([
					attribute.src(image),
					attribute.class("h-32 rounded-md")
				])
			},
			html.p([], [element.text(text)])
		])
	])
}

fn button(text: String, color: String, link: String) -> element.Element(a) {
	html.a([attribute.href(link), attribute.target("_blank")], [
		html.button([attribute.class(color <> " rounded-md backdrop-blur-md py-2 px-4 font-semibold")], [element.text(text)])
	])
}
