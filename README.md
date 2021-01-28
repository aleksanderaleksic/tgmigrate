# tgmigrate

tgmigrate helps you with migrating state in and across terraform state files in a terragrunt environment.</br>

#### Why?

Migrating state is necessary when refactoring your terragrunt project, and moving resources between state files or
within them can be tricky. There are projects like [tfmigrate](https://github.com/minamijoyo/tfmigrate) that's targeted
towards terraform, but we found it hard to use with terragrunt. Usually we write shell scripts to move or remove state
resources with the `terraform state` commands, but there is no way to store the history to check if the migration have
been applied.

## Features

- Simple and declarative migrations
- Migration history
- Automation friendly
- AWS s3 support
- Test the migration with dryrun the option

## Install

Download the latest package from [releases](https://github.com/aleksanderaleksic/tgmigrate/releases) and put it in the
executable path.

## Usage

### You first have to set up a `.tgmigrate.hcl` file:

Example config file:

```hcl
migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "airthings-terraform-states-${ACCOUNT}"
      region = "us-east-1"
      assume_role = "${ASSUME_ROLE}"
      key = "history.json"
    }
  }

  state "s3" {
    bucket = "airthings-terraform-states-${ACCOUNT}"
    region = "us-east-1"
    assume_role = "${ASSUME_ROLE}"
  }
}
```

this file can be specified by using the `-c`(config) flag, </br> if not specified tgmigrate will look ub the parent
folders up to `$HOME`

`migration = "./migration"` refers to the directory where the migration files are located, this path is relative to the
config file. </br>
`history` is where you can configure where to store the history, currently its only `storage "s3"` that's
supported.</br>
`state` is where you configure where your state is located, also here its only `state "s3"` that's supported.

A nice feature in the config file is the support for variables using the common `${variable_name}` syntax, How to
populate the variables is described further down.

### Next you need to make a migration file:

Under the migration directory you specified in the config file, create a new file and name it something that makes sense
for you. I like to use today's data date + a short description, this makes sure the migration's comes in order by data,
but it doesn't really matter.

Example migration file:

```hcl
migration {
  environments = [
    "dev",
    "prod",
    "developer"
  ]
  description = <<EOF
    - Move the rest2 lambda to a separate module.
EOF
}

migrate "move" "rest_2" {
  from {
    state = "us-east-1/apis/rest"
    resource = "aws_lambda_function.rest_2_api_lambda"
  }

  to {
    state = "us-east-1/apis/rest_2"
    resource = "aws_lambda_function.rest_api_lambda"
  }
}

migrate "remove" "api_gateway_integration_for_rest_2" {
  state = "us-east-1/apis/rest"
  resource = "aws_apigatewayv2_integration.rest_2_api"
}
```

The `migration` block contains metadata for the migration file. Here you can specify the `environments` the migration
file is intended for, this is optional to include. You can also include a description for the migration but that is
optional.

The `migrate` blocks have 2 supported types: `move` and `remove`.</br>
On `move` you need to specify a `from` and `to` block, while on `remove` you can provide the `state` and `resource`
variable. </br>
In both cases `state` refers to the location within the s3_bucket. `resource` refers to the resource type + name in the
state file.

### Integrate with terragrunt:

tgmigrate integrates well with terragrunt by using the before hook.</br>

```hcl
before_hook "dryrun_migrations" {
  commands = [
    "plan"
  ]
  execute = [
    "tgmigrate",
    "-d",
    "-y",
    "--cv=ACCOUNT_ID=${local.account_id};ASSUME_ROLE=${local.terraform_role_arn}",
    "apply",
    "prod"
  ]
  run_on_error = false
}
before_hook "run_migrations" {
  commands = [
    "plan"
  ]
  execute = [
    "tgmigrate",
    "-y",
    "--cv=ACCOUNT_ID=${local.account_id};ASSUME_ROLE=${local.terraform_role_arn}",
    "apply",
    "prod"
  ]
  run_on_error = false
}
```

Notice the `--cv`(config-variables) flag, here we specify the `ACCOUNT_ID` and the `ASSUME_ROLE` variable that we use in
the config file.</br>
Also notice the `prod` sub-command, this is telling tgmigrate to only apply migrations for the prod environment. </br>
This means that the migration file above would be applied in this case, because it has `prod` in the environments list.