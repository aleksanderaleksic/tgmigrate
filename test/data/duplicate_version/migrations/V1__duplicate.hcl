migration {
  environments = [
    "test"
  ]
  description = <<EOF
    - Remove testfile from files module
EOF
}

migrate "remove" "file" {
  state = "us-east-1/files"
  resource = "file.test_file"
}