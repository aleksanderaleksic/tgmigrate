terraform {
  backend "local" {}
}

resource "local_file" "test_file1" {
  filename = "test_file1.txt"
  content = "hello"
}

resource "local_file" "test_file2" {
  filename = "test_file2.txt"
  content = "hello"
}