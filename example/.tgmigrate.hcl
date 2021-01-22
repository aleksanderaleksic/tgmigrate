migration {
  migration = "./migrations"

  history {
    storage "local" {
      path = "history.json"
    }
  }

  state "s3" {
    bucket = "airthings-terraform-states-512741945438"
    region = "us-east-1"
    assume_role = "arn:aws:iam::512741945438:role/admin"
  }
}