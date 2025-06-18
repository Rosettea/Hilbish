import gleam/int
import gleam/option
import gleam/string

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
