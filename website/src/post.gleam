import glaml
import gleam/option

pub type Post {
  Post(
    name: String,
    title: String,
    slug: String,
    metadata: option.Option(glaml.Document),
    contents: String,
  )
}
