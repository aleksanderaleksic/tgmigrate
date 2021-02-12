migration {
  environments = [
    "run2"
  ]
  description = <<EOF
    - Revert v1 migration
EOF
}

migrate "move" "file_from_file2_to_file1" {
  from {
    state = "file2"
    resource = "local_file.test_file"
  }

  to {
    state = "file1"
    resource = "local_file.test_file"
  }
}