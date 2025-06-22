import gleam/int
import gleam/option
import gleam/string

import lustre/attribute
import lustre/element
import lustre/element/html

import glaml
import post

pub fn sort_weight(p1: #(String, post.Post), p2: #(String, post.Post)) {
  let assert option.Some(p1_metadata) = { p1.1 }.metadata
  let p1_weight = case
    glaml.select_sugar(glaml.document_root(p1_metadata), "weight")
  {
    Ok(glaml.NodeInt(w)) -> w
    _ -> 0
  }

  let assert option.Some(p2_metadata) = { p2.1 }.metadata
  let p2_weight = case
    glaml.select_sugar(glaml.document_root(p2_metadata), "weight")
  {
    Ok(glaml.NodeInt(w)) -> w
    _ -> 0
  }

  case p1_weight == p2_weight {
    True -> string.compare({ p1.1 }.name, { p2.1 }.name)
    False -> int.compare(p1_weight, p2_weight)
  }
}

pub fn link(url: String, text: String, out: Bool) {
  html.a(
    [
      attribute.href(url),
      case out {
        False -> attribute.none()
        True -> attribute.target("_blank")
      },
    ],
    [
      html.span(
        [
          attribute.class(
            "inline-flex text-light dark:text-pink-300 dark:hover:text-pink-200 text-pink-600 hover:text-pink-500 hover:underline",
          ),
        ],
        [
          case out {
            False -> element.none()
            True ->
              element.unsafe_raw_html(
                "",
                "tag",
                [],
                "<svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\">
  <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25\" />
</svg>
",
              )
          },
          element.text(text),
        ],
      ),
    ],
  )
}
