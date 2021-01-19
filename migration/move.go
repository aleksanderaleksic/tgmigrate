package migration

import "github.com/hashicorp/hcl/v2/gohcl"

type MoveBlock struct {
	From MoveFromBlock `hcl:"from,block"`
	To   MoveFromBlock `hcl:"to,block"`
}

type MoveFromBlock struct {
	State    string `hcl:"state"`
	Resource string `hcl:"resource"`
}

type MoveToBlock struct {
	State    string `hcl:"state"`
	Resource string `hcl:"resource"`
}

func ParseMigrateMoveBlock(block MigrationsBlock) (*MoveBlock, error) {
	var config MoveBlock
	diags := gohcl.DecodeBody(block.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}
