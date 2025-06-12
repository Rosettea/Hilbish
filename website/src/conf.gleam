pub const base_url = "http://localhost:9080"

pub fn base_url_join(cont: String) -> String {
  base_url <> "/" <> cont
}
