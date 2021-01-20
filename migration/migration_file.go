package migration

import (
	"crypto/sha256"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func GetMigrationFiles(cfg config.Config) (*[]File, error) {
	var migrationFiles []File

	err := filepath.Walk(cfg.AbsoluteMigrationDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, "hcl") {
				return nil
			}

			source, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			migrationFile, err := parseMigrationFile(info.Name(), source)
			if err != nil {
				return err
			}

			migrationFiles = append(migrationFiles, *migrationFile)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return &migrationFiles, err
}

func parseMigrationFile(filename string, source []byte) (*File, error) {
	fileSha256 := fmt.Sprintf("%x", sha256.Sum256(source))

	var f FileContent
	err := hclsimple.Decode(filename, source, nil, &f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode migration file: %s, err: %s", filename, err)
	}

	var fileMeta = FileMetadata{
		FileName: filename,
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
