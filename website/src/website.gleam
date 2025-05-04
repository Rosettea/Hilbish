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
	html.html([attribute.class("bg-stone-50 dark:bg-stone-950 text-black dark:text-white")], [
		html.head([], [
			html.meta([
				attribute.name("viewport"),
				attribute.attribute("content", "width=device-width, initial-scale=1.0")
			]),
			html.link([
				attribute.rel("stylesheet"),
				attribute.href("./tailwind.css")
			]),
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
