import gleam/option
import gleam/dict
import gleam/io
import gleam/order
import gleam/list
import gleam/string

import lustre/attribute
import lustre/element
import lustre/element/html
import lustre/ssg
import lustre/ssg/djot
import tom
import simplifile
import glaml

import post
import pages/index
import pages/doc

const base_url = "https://rosettea.github.io/Hilbish/versions/new-website"

pub fn main() {
	let assert Ok(files) = simplifile.get_files("./content")
	let posts = list.map(files, fn(path: String) {
		let assert Ok(ext) = path |> string.split(".") |> list.last
		let slug = path |> string.replace("./content", "") |> string.drop_end({ext |> string.length()} + 1)
		let assert Ok(name) = slug |> string.split("/") |> list.last

		let assert Ok(content) = simplifile.read(path)
		let frontmatter = djot.frontmatter(content)
		let metadata = case frontmatter {
			Ok(frntmtr) -> {
				let assert Ok([metadata]) = glaml.parse_string(frntmtr)
				option.Some(metadata)
			}
			Error(_) -> option.None
		}
		let content = djot.content(content)

		let title = case metadata {
			option.Some(metadata) -> {
				case glaml.select_sugar(glaml.document_root(metadata), "title") {
					Ok(glaml.NodeStr(s)) -> s
					_ -> ""
				}
				
			}
			option.None -> ""
		}

		let assert Ok(filename) = path |> string.split("/") |> list.last
		#(slug, post.Post(name, title, slug, metadata, content))
	})

	let doc_pages = list.filter(posts, fn(page) {
		let isdoc = is_doc_page(page.0)
		io.debug(page.0)
		io.debug(isdoc)
		isdoc
	}) |> list.filter(fn(page) {
		case page.1.metadata {
			option.Some(_) -> True
			option.None -> False
		}
	}) |> list.sort(fn(p1, p2) {
		io.debug(p1)
		io.debug(p2)
		let assert option.Some(p1_metadata) = p1.1.metadata
		let p1_weight = case glaml.select_sugar(glaml.document_root(p1_metadata), "weight") {
			Ok(glaml.NodeInt(w)) -> w
			_ -> 0
		}

		let assert option.Some(p2_metadata) = p2.1.metadata
		let p2_weight = case glaml.select_sugar(glaml.document_root(p2_metadata), "weight") {
			Ok(glaml.NodeInt(w)) -> w
			_ -> 0
		}

		case p1_weight == 0 {
			True -> order.Eq
			False -> {
				case p1_weight < p2_weight {
					True -> order.Lt
					False -> order.Gt
				}
			}
		}
	})

	let build = ssg.new("./public")
	|> ssg.add_static_dir("static")
	|> ssg.add_static_route("/", create_page(index.page()))
	|> list.fold(posts, _, fn(config, post) {
		let route = case post.1.name {
			"_index" -> post.0 |> string.drop_end("_index" |> string.length())
			_ -> post.0
		}
		

		let page = case is_doc_page(post.0) {
			True -> doc.page(post.1, doc_pages)
			False -> doc.page(post.1, doc_pages)
		}
		ssg.add_static_route(config, route, create_page(page))
	})
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

fn is_doc_page(slug: String) {
	let is_docs = case slug {
		"/docs" <> _ -> True
		_ -> False
	}
}

fn base_url_join(cont: String) -> String {
	return base_url <> "/" <> cont
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
				attribute.href(base_url_join("tailwind.css"))
			]),
			html.title([], "Hilbish"),
			html.meta([attribute.name("theme-color"), attribute.content("#ff89dd")]),
			html.meta([attribute.content(base_url_join("hilbish-flower.png")), attribute.attribute("property", "og:image")]),
			html.meta([attribute.content("Hilbish"), attribute.attribute("property", "og:title")]), // this should be same as title
			html.meta([attribute.content("Hilbish"), attribute.attribute("property", "og:site_name")]),
			html.meta([attribute.content("website"), attribute.attribute("property", "og:type")]),
			html.meta([attribute.content(description), attribute.attribute("property", "og:description")]),
			html.meta([attribute.content(description), attribute.name("description")]),
			html.meta([attribute.name("keywords"), attribute.content("Lua,Shell,Hilbish,Linux,zsh,bash")]),
			html.meta([attribute.content(base_url), attribute.attribute("property", "og:url")])
		]),
		html.body([], [
			html.nav([attribute.class("flex sticky top-0 w-full z-50 border-b border-b-zinc-300 backdrop-blur-md h-12")], [
				html.div([attribute.class("flex my-auto px-2")], [
					html.div([], [
						html.a([attribute.href("/"), attribute.class("flex items-center gap-1")], [
							html.img([
								attribute.src("/hilbish-flower.png"),
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
			content,
			html.footer([attribute.class("py-4 px-6 flex flex-row justify-around border-t border-t-zinc-300")], [
				html.div([attribute.class("flex flex-col")], [
					html.a([attribute.href(base_url), attribute.class("flex items-center gap-1")], [
						html.img([
							attribute.src("/hilbish-flower.png"),
							attribute.class("h-24")
						]),
						html.span([
							attribute.class("self-center text-6xl")
						], [
							element.text("Hilbish"),
						]),
					]),
					html.span([attribute.class("text-xl")], [element.text("The Moon-powered shell!")]),
					html.span([attribute.class("text-light text-neutral-500")], [element.text("MIT License, copyright sammyette 2025")])
				]),
				html.div([attribute.class("flex flex-col")], [
					link("https://github.com/Rosettea/Hilbish", "GitHub")
				])
			])
		])
	])
}

fn link(url: String, text: String) {
	html.a([attribute.href(url)], [
		html.span([attribute.class("text-pink-300 text-light")], [element.text(text)])
	])
}
