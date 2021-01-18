migration_config {
  environments = [
    "dev",
    "prod",
    "developer"
  ]
  description = <<EOF
    - Move test.com to global dns state,
    - Remove example.com from us-east-1
    - Import something.com to global
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
  resources = "module.example_com"
}

data "aws_route53_zone" "something_com" {
  name = "something.com"
}

migrate "import" "something_com" {
  state = "global/dns"
  resource = "module.something_com.aws_route53_zone.something_com"
  import_value = data.aws_route53_zone.something_com.id
}
