terraform {
  backend "local" {}
}

resource "local_file" "test_file" {
  filename = "test_file.txt"
  content = "hello"
}