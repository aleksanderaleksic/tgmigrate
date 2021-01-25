module github.com/aleksanderaleksic/tgmigrate

go 1.15

require (
	github.com/aws/aws-sdk-go v1.35.28
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/hashicorp/terraform-exec v0.12.0
	github.com/seqsense/s3sync v1.8.0
	github.com/urfave/cli/v2 v2.3.0
	github.com/zclconf/go-cty v1.2.1
)

replace github.com/hashicorp/terraform-exec => github.com/aleksanderaleksic/terraform-exec v0.12.0