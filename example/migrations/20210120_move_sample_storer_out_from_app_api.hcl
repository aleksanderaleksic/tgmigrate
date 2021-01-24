migration {
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

migrate "move" "store_samples_lambda" {
  from {
    state = "us-east-1/apis/app"
    resource = "aws_lambda_function.store_samples_lambda"
  }

  to {
    state = "us-east-1/infrastructure/store-samples"
    resource = "aws_lambda_function.store_samples_lambda"
  }
}

migrate "move" "store_samples_cloudwatch_group" {
  from {
    state = "us-east-1/apis/app"
    resource = "aws_cloudwatch_log_group.store_samples_cloudwatch_group"
  }

  to {
    state = "us-east-1/infrastructure/store-samples"
    resource = "aws_cloudwatch_log_group.store_samples_cloudwatch_group"
  }
}

migrate "move" "store_samples_lambda_log_filter" {
  from {
    state = "us-east-1/apis/app"
    resource = "aws_cloudwatch_log_subscription_filter.store_samples_lambda_log_filter"
  }

  to {
    state = "us-east-1/infrastructure/store-samples"
    resource = "aws_cloudwatch_log_subscription_filter.store_samples_lambda_log_filter"
  }
}

migrate "move" "store_samples_lambda_alias" {
  from {
    state = "us-east-1/apis/app"
    resource = "aws_lambda_alias.store_samples_lambda_alias"
  }

  to {
    state = "us-east-1/infrastructure/store-samples"
    resource = "aws_lambda_alias.store_samples_lambda_alias"
  }
}