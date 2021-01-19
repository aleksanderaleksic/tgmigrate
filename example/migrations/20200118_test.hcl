migration_config {
  environments = [
    "dev",
    "prod",
    "developer"
  ]
  description = <<EOF
    - Move test.com to global dns state,
    - Remove example.com from us-east-1
EOF
}

migrate "move" "test_com" {
  from {
    state = "us-east-1/dns"
    resource = "module.test_com.aws_route53_zone.test_com"
  }

  to {
    state = "global/dns"
    resource = "module.test_com.aws_route53_zone.test_com"
  }
}

migrate "remove" "example_com" {
  state = "us-east-1/dns"
  resource = "module.example_com"
}
