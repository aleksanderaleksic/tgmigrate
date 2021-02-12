migration {
  environments = [
    "run1"
  ]
  description = <<EOF
    - Remove testfile from files module
EOF
}

migrate "remove" "file" {
  state = "file3"
  resource = "local_file.test_file1"
}