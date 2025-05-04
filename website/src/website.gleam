import gleam/io

import lustre/attribute
import lustre/element
import lustre/element/html
import lustre/ssg

import pages/index

pub fn main() {
	let build = ssg.new("./public")
	|> ssg.add_static_dir("static")
	|> ssg.add_static_route("/", create_page(index.page()))
	|> ssg.use_index_routes
	|> ssg.build

	case build {
		Ok(_) -> io.println("Website successfully built!")
		Error(e) -> {
			io.debug(e)
			io.println("Website could not be built.")
		}
	}
}

fn create_page(content: element.Element(a)) -> element.Element(a) {
	let description = "Something Unique. Hilbish is the new interactive shell for Lua fans. Extensible, scriptable, configurable: All in Lua."

	html.html([attribute.class("bg-stone-50 dark:bg-neutral-950 text-black dark:text-white")], [
		html.head([], [
			html.meta([
				attribute.name("viewport"),
				attribute.attribute("content", "width=device-width, initial-scale=1.0")
			]),
			html.link([
				attribute.rel("stylesheet"),
				attribute.href("./tailwind.css")
			]),
			html.title([], "Hilbish"),
			html.meta([attribute.name("theme-color"), attribute.content("#ff89dd")]),
			html.meta([attribute.content("./hilbish-flower.png"), attribute.attribute("property", "og:image")]),
			html.meta([attribute.content("Hilbish"), attribute.attribute("property", "og:title")]), // this should be same as title
			html.meta([attribute.content("Hilbish"), attribute.attribute("property", "og:site_name")]),
			html.meta([attribute.content("website"), attribute.attribute("property", "og:type")]),
			html.meta([attribute.content(description), attribute.attribute("property", "og:description")]),
			html.meta([attribute.content(description), attribute.name("description")]),
			html.meta([attribute.name("keywords"), attribute.content("Lua,Shell,Hilbish,Linux,zsh,bash")]),
			html.meta([attribute.content("https://rosettea.github.io/Hilbish/versions/new-website"), attribute.attribute("property", "og:url")])
		]),
		html.body([], [
			html.nav([attribute.class("fixed top-0 w-full z-50 p-1 mb-2 border-b border-b-zinc-300 backdrop-blur-md")], [
				html.div([attribute.class("flex mx-auto")], [
					html.div([], [
						html.a([attribute.href("/"), attribute.class("flex items-center gap-1")], [
							html.img([
								attribute.src("./hilbish-flower.png"),
								attribute.class("h-6")
							]),
							html.span([
								attribute.class("self-center text-2xl")
							], [
								element.text("Hilbish"),
							]),
						]),
					])
				]),
			]),
			html.main([attribute.class("mx-4")], [content])
		])
	])
}
