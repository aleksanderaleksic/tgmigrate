package migration

import (
	"crypto/sha256"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type FileContent struct {
	Config     Config            `hcl:"migration_config,block"`
	Migrations []MigrationsBlock `hcl:"migrate,block"`
}

type File struct {
	Metadata   FileMetadata
	Config     Config
	Migrations []Migration
}

func ParseMigrationFile(filename string, source []byte) (*File, error) {
	fileSha256 := fmt.Sprintf("%x", sha256.Sum256(source))

	var f FileContent
	err := hclsimple.Decode(filename, source, nil, &f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode migration file: %s, err: %s", filename, err)
	}

	var fileMeta = FileMetadata{
		FileHash: fileSha256,
	}

	var migrations []Migration

	for _, migration := range f.Migrations {
		switch migration.Type {
		case "move":
			move, err := ParseMigrateMoveBlock(migration)
			if err != nil {
				return nil, err
			}
			var m = Migration{
				Type:   migration.Type,
				Name:   migration.Name,
				Move:   move,
				Remove: nil,
			}
			migrations = append(migrations, m)

		case "remove":
			remove, err := ParseMigrateRemoveBlock(migration)
			if err != nil {
				return nil, err
			}
			var m = Migration{
				Type:   migration.Type,
				Name:   migration.Name,
				Move:   nil,
				Remove: remove,
			}
			migrations = append(migrations, m)
		}
	}

	file := File{
		Metadata:   fileMeta,
		Config:     f.Config,
		Migrations: migrations,
	}

	return &file, nil
}
