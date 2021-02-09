terraform {
  backend "local" {}
}

resource "local_file" "file_map" {
  for_each = toset([
    "file1",
    "file2",
    "file3"])
  filename = "test_${each.key}.txt"
  content = each.key
}