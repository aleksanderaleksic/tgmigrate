migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "airthings-terraform-states-${ACCOUNT}"
      region = "us-east-1"
      assume_role = "arn:aws:iam::${ACCOUNT}:role/admin"
      key = "history.json"
    }
  }

  state "s3" {
    bucket = "airthings-terraform-states-${ACCOUNT}"
    region = "us-east-1"
    assume_role = "arn:aws:iam::${ACCOUNT}:role/admin"
  }
}