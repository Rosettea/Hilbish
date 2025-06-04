import gleam/dict
import gleam/list
import gleam/string

import lustre/attribute
import lustre/element
import lustre/element/html
import lustre/ssg/djot

import jot
import post

pub fn page(p: post.Post, doc_pages_list) -> element.Element(a) {
	html.div([attribute.class("flex flex-col")], [
		html.div([attribute.class("h-10 flex py-2 px-4 border-b border-b-zinc-300 w-full gap-2 backdrop-blur-sm bg-zinc-300/50 dark:bg-zinc-800/50 z-50")], [
			html.label([attribute.for("sidebar-toggle"), attribute.class("cursor-pointer")], [
				element.unsafe_raw_html("", "tag", [], "<svg xmlns=\"http://www.w3.org/2000/svg\" height=\"24px\" viewBox=\"0 -960 960 960\" width=\"24px\" class=\"fill-black\"><path d=\"M120-240v-80h240v80H120Zm0-200v-80h480v80H120Zm0-200v-80h720v80H120Z\"/></svg>"),
			]),
			html.span([], [element.text(p.title)])
		]),
		html.div([attribute.class("grid")], [
			html.input([attribute.type_("checkbox"), attribute.id("sidebar-toggle"), attribute.class("peer hidden")]),
			html.div([attribute.class("border-r border-r-zinc-300 col-start-1 row-start-1 sticky top-22 sm:top-12 h-full sm:h-svh bg-neutral-200 dark:bg-neutral-900 basis-3/5 transition-transform duration-300 -translate-x-full peer-checked:translate-x-0 z-30")], [
				html.div([attribute.class("p-4 -mb-4 overflow-y-auto h-full")], [
					html.h2([attribute.class("text-xl font-semibold mb-4")], [element.text("Sidebar")]),
					html.ul([], list.map(doc_pages_list, fn(post: #(String, post.Post)) {
						html.li([attribute.class("mb-2")], [element.text(post.1.title)])
					}))
				])
			]),
			html.main([attribute.class("col-start-1 row-start-1 transition-all duration-300 peer-checked:filter peer-checked:blur-sm peer-checked:bg-black/30 px-4 pt-2")], [
				html.h1([attribute.class("font-bold text-4xl")], [element.text(p.title)]),
				// todo: add date of publishing
				//html.time([], [])
				//html.small([], [element.text({{p.contents |> string.split(" ") |> list.length} / 200} |> int.to_string <> " min read")]),
				//element.unsafe_raw_html("namespace", "Tag", [], render_doc(p.contents))
				..render_doc(p.contents)
			])
		])
	])
}

fn render_doc(md: String) {
	let renderer = djot.Renderer(
		..djot.default_renderer(),
		heading: fn(attrs, level, content) {
			let size = case level {
				1 -> "text-4xl"
				2 -> "text-3xl"
				3 -> "text-2xl"
				_ -> "text-xl"
			}
			let attr = dict.insert(attrs, "class", "font-bold " <> size)

			case level {
				1 -> html.h1(to_attr(attr), content)
				2 -> html.h2(to_attr(attr), content)
				3 -> html.h3(to_attr(attr), content)
				_ -> html.p(to_attr(attr), content)
			}
		}
	)
	djot.render(md, renderer)
}

fn to_attr(attrs) {
	use attrs, key, val <- dict.fold(attrs, [])
	[attribute.attribute(key, val), ..attrs]
}

fn render_doc_(md: String) -> String {
	// w-full m-2 p-2 bg-neutral-700
	let doc = jot.parse(md)
	let updated_content = list.map(doc.content, fn(container) {
		case container {
			jot.Heading(attributes, level, content) -> {
				let size = case level {
					1 -> "text-4xl"
					2 -> "text-3xl"
					3 -> "text-2xl"
					_ -> "text-xl"
				}
				let attr = dict.insert(attributes, "class", "font-bold " <> size)
				jot.Heading(attr, level, content)
			}
			_ -> container
		}
	})
	echo doc

	jot.document_to_html(jot.Document(
		content: updated_content,
		references: doc.references,
		footnotes: doc.footnotes
	))
}

fn page_(p: post.Post, doc_pages_list) -> element.Element(a) {
	html.div([attribute.class("relative h-screen flex")], [
		html.div([attribute.class("-mt-2 -mx-4 py-2 px-4 border-b border-b-zinc-300 flex gap-2 font-semibold")], [
			html.label([attribute.for("sidebar-toggle"), attribute.class("cursor-pointer")], [
				element.unsafe_raw_html("", "tag", [], "<svg xmlns=\"http://www.w3.org/2000/svg\" height=\"24px\" viewBox=\"0 -960 960 960\" width=\"24px\" class=\"fill-black\"><path d=\"M120-240v-80h240v80H120Zm0-200v-80h480v80H120Zm0-200v-80h720v80H120Z\"/></svg>"),
			]),
			html.span([], [element.text(p.title)])
		]),
		html.div([attribute.class("relative flex")], [
			html.div([attribute.class("absolute top-0 left-0 h-full bg-gray-200 w-64 transition-transform duration-300 -translate-x-full peer-checked:translate-x-0 z-30")], [
				html.div([attribute.class("p-4")], [
					html.h2([attribute.class("text-xl font-semibold mb-4")], [element.text("Sidebar")]),
					html.ul([], [
						html.li([attribute.class("mb-2")], [element.text("Test")])
					])
				])
			]),
			html.input([attribute.type_("checkbox"), attribute.id("sidebar-toggle"), attribute.class("peer hidden")]),
			html.main([attribute.class("flex-1 transition-all duration-300 peer-checked:filter peer-checked:blur-sm peer-checked:opacity-50")], [
				html.h1([], [element.text(p.title)]),
				// todo: add date of publishing
				//html.time([], [])
				//html.small([], [element.text({{p.contents |> string.split(" ") |> list.length} / 200} |> int.to_string <> " min read")]),
				//element.unsafe_raw_html("namespace", "Tag", [], md.md_to_html(p.contents))
			])
		])
	])
}
