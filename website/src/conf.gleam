pub const base_url = "https://rosettea.github.io/Hilbish/versions/new-website/"

pub fn base_url_join(cont: String) -> String {
  base_url <> "/" <> cont
}
