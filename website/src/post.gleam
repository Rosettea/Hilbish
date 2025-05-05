import gleam/option
import glaml

pub type Post {
	Post(name: String, title: String, slug: String, metadata: option.Option(glaml.Document), contents: String)
}
