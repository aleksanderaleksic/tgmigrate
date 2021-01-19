package migration

type Config struct {
	Environments []string `hcl:"environments,optional"`
	Description  string   `hcl:"description,optional"`
}
