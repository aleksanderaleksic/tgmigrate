migration {
  migration_dir = "migrations"

  history {
    storage "local" {
      path = "history.json"
    }
  }

  state "local" {
    directory = "state"
    state_file_name = "terraform.tfstate"
  }
}