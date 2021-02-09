migration {
  environments = [
    "test"
  ]
  description = <<EOF
    - Move file between modules
EOF
}

migrate "move" "file_from_file1_to_file2" {
  from {
    state = "file1"
    resource = "local_file.test_file"
  }

  to {
    state = "file2"
    resource = "local_file.test_file"
  }
}