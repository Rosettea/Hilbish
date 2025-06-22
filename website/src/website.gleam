import gleam/io
import gleam/list
import gleam/option
import gleam/order
import gleam/string
import util

import glaml
import lustre/attribute
import lustre/element
import lustre/element/html
import lustre/ssg
import lustre/ssg/djot
import simplifile
import tom

import conf
import pages/doc
import pages/index
import post

pub fn main() {
  let assert Ok(files) = simplifile.get_files("./content")
  let posts =
    list.map(files, fn(path: String) {
      let assert Ok(ext) = path |> string.split(".") |> list.last
      let slug =
        path
        |> string.replace("./content", "")
        |> string.drop_end({ ext |> string.length() } + 1)
      let assert Ok(name) = slug |> string.split("/") |> list.last

      let slug = case name {
        "_index" -> slug |> string.drop_end({ "_index" |> string.length() } + 1)
        _ -> slug
      }

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

      let description = case metadata {
        option.Some(metadata) -> {
          case
            glaml.select_sugar(glaml.document_root(metadata), "description")
          {
            Ok(glaml.NodeStr(s)) -> s
            _ -> ""
          }
        }
        option.None -> ""
      }

      let assert Ok(filename) = path |> string.split("/") |> list.last
      #(slug, post.Post(name, description, title, slug, metadata, content))
    })

  let doc_pages =
    list.filter(posts, fn(page) {
      let isdoc = is_doc_page(page.0)
      //io.debug(page.0)
      //io.debug(isdoc)
      isdoc
    })
    |> list.filter(fn(page) {
      case { page.1 }.metadata {
        option.Some(_) -> True
        option.None -> False
      }
    })
    |> list.sort(util.sort_weight)

  let build =
    ssg.new("./public")
    |> ssg.add_static_dir("static")
    |> ssg.add_static_route("/", create_page(index.page(), False))
    |> list.fold(posts, _, fn(config, post) {
      let page = case is_doc_page(post.0) {
        True -> doc.page(post.1, post.0, doc_pages)
        False -> doc.page(post.1, post.0, doc_pages)
      }
      //io.debug(post.0)
      ssg.add_static_route(
        config,
        post.0,
        create_page(page, is_doc_page(post.0)),
      )
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
  case slug {
    "/docs" <> _ -> True
    _ -> False
  }
}

fn create_page(
  content: element.Element(a),
  doc_page: Bool,
) -> element.Element(a) {
  let description =
    "Something Unique. Hilbish is the new interactive shell for Lua fans. Extensible, scriptable, configurable: All in Lua."

  html.html(
    [
      attribute.class(
        "bg-stone-50 dark:bg-neutral-900 text-black dark:text-white",
      ),
    ],
    [
      html.head([], [
        html.meta([
          attribute.name("viewport"),
          attribute.attribute(
            "content",
            "width=device-width, initial-scale=1.0",
          ),
        ]),
        html.link([
          attribute.rel("stylesheet"),
          attribute.href(conf.base_url_join("/tailwind.css")),
        ]),
        html.title([], "Hilbish"),
        html.meta([attribute.name("theme-color"), attribute.content("#ff89dd")]),
        html.meta([
          attribute.content(conf.base_url_join("/hilbish-flower.png")),
          attribute.attribute("property", "og:image"),
        ]),
        html.meta([
          attribute.content("Hilbish"),
          attribute.attribute("property", "og:title"),
        ]),
        // this should be same as title
        html.meta([
          attribute.content("Hilbish"),
          attribute.attribute("property", "og:site_name"),
        ]),
        html.meta([
          attribute.content("website"),
          attribute.attribute("property", "og:type"),
        ]),
        html.meta([
          attribute.content(description),
          attribute.attribute("property", "og:description"),
        ]),
        html.meta([
          attribute.content(description),
          attribute.name("description"),
        ]),
        html.meta([
          attribute.name("keywords"),
          attribute.content("Lua,Shell,Hilbish,Linux,zsh,bash"),
        ]),
        html.meta([
          attribute.content(conf.base_url),
          attribute.attribute("property", "og:url"),
        ]),
      ]),
      html.body([attribute.class("h-screen flex flex-col")], [
        util.nav(),
        content,
        case doc_page {
          True -> element.none()
          False -> util.footer()
        },
      ]),
    ],
  )
}
