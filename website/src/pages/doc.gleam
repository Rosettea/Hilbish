import gleam/list
import gleam/string

import lustre/attribute
import lustre/element
import lustre/element/html

import md
import post

pub fn page(p: post.Post, doc_pages_list) -> element.Element(a) {
	html.div([attribute.class("flex flex-col")], [
		html.div([attribute.class("block sm:hidden h-10 sticky top-12 flex py-2 px-4 border-b border-b-zinc-300 w-full gap-2 backdrop-blur-sm bg-zinc-300/50 z-50")], [
			html.label([attribute.for("sidebar-toggle"), attribute.class("cursor-pointer")], [
				element.unsafe_raw_html("", "tag", [], "<svg xmlns=\"http://www.w3.org/2000/svg\" height=\"24px\" viewBox=\"0 -960 960 960\" width=\"24px\" class=\"fill-black\"><path d=\"M120-240v-80h240v80H120Zm0-200v-80h480v80H120Zm0-200v-80h720v80H120Z\"/></svg>"),
			]),
			html.span([], [element.text(p.title)])
		]),
		html.div([attribute.class("grid sm:flex")], [
			html.input([attribute.type_("checkbox"), attribute.id("sidebar-toggle"), attribute.class("peer hidden")]),
			html.div([attribute.class("border-r border-r-zinc-300 col-start-1 row-start-1 sticky top-22 sm:top-12 h-full sm:h-svh bg-neutral-200 basis-3/5 transition-transform duration-300 -translate-x-full sm:translate-x-0 peer-checked:translate-x-0 z-30")], [
				html.div([attribute.class("p-4 -mb-4 overflow-y-auto h-full")], [
					html.h2([attribute.class("text-xl font-semibold mb-4")], [element.text("Sidebar")]),
					html.ul([], list.map(doc_pages_list, fn(post: #(String, post.Post)) {
						html.li([attribute.class("mb-2")], [element.text(post.1.title)])
					}))
				])
			]),
			html.main([attribute.class("col-start-1 row-start-1 transition-all duration-300 peer-checked:filter peer-checked:blur-sm peer-checked:bg-black/30 px-4 pt-2")], [
				html.h1([], [element.text(p.title)]),
				// todo: add date of publishing
				//html.time([], [])
				//html.small([], [element.text({{p.contents |> string.split(" ") |> list.length} / 200} |> int.to_string <> " min read")]),
				element.unsafe_raw_html("namespace", "Tag", [], md.md_to_html(p.contents))
			])
		])
	])
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
