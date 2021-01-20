migration_config {
  environments = [
    "dev",
    "prod",
    "developer"
  ]
  description = <<EOF
    - Move test_file_1 to global dns state,
    - Remove test_file_123 form us-east-1
EOF
}

migrate "move" "test_file_1" {
  from {
    state = "us-east-1/file"
    resource = "local_file.test"
  }

  to {
    state = "global/file"
    resource = "local_file.test"
  }
}

migrate "remove" "test_file_123" {
  state = "us-east-1/file"
  resource = "local_file.test123"
}
