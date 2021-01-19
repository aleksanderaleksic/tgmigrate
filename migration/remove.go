package migration

import "github.com/hashicorp/hcl/v2/gohcl"

type RemoveBlock struct {
	State    string `hcl:"state"`
	Resource string `hcl:"resource"`
}

func ParseMigrateRemoveBlock(block MigrationsBlock) (*RemoveBlock, error) {
	var config RemoveBlock
	diags := gohcl.DecodeBody(block.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}
