migration {
  environments = [
    "test"
  ]
  description = <<EOF
    - Move testfile from files module
EOF
}

migrate "remove" "file" {
  state = "us-east-1/files"
}