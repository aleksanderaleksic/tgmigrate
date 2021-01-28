migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "airthings-terraform-states-${ACCOUNT_ID}"
      region = "us-east-1"
      assume_role = "arn:aws:iam::${ACCOUNT_ID}:role/admin"
      key = "history.json"
    }
  }

  state "s3" {
    bucket = "airthings-terraform-states-${ACCOUNT_ID}"
    region = "us-east-1"
    assume_role = "arn:aws:iam::${ACCOUNT_ID}:role/admin"
  }
}