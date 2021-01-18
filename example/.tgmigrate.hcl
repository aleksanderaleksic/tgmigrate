migration {
  migration_dir = "./migrations"

  history {
    storage "s3" {
      key = "migrations/history.json"
    }
  }
}