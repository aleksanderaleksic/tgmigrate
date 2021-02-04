migration {
  environments = [
    "test"
  ]
  description = <<EOF
    - Move rest_api lambda to rest_v2
EOF
}

migrate "move" "rest_api" {
  to {
    state = "us-east-1/apis/rest_v2"
    resource = "aws_lambda_function.rest_api"
  }
}