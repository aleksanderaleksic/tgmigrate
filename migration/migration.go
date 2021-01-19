package migration

import "github.com/hashicorp/hcl/v2"

type MigrationsBlock struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Remain hcl.Body `hcl:",remain"`
}

type Migration struct {
	Type   string
	Name   string
	Move   *MoveBlock
	Remove *RemoveBlock
}
