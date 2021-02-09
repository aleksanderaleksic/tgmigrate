migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "tgmigrate-e2e-test-bucket"
      region = "eu-central-1"
      key = "${TEST_ID}/history.json"
    }
  }

  state "s3" {
    bucket = "tgmigrate-e2e-test-bucket"
    prefix = "${TEST_ID}"
    region = "eu-central-1"
  }
}