package migration

import (
	"crypto/sha256"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type FileContent struct {
	Config     Config            `hcl:"migration,block"`
	Migrations []MigrationsBlock `hcl:"migrate,block"`
}

type File struct {
	Metadata   FileMetadata
	Config     Config
	Migrations []Migration
}

//Sort interface implementation
type FilesBySequence []File

func (f FilesBySequence) Len() int           { return len(f) }
func (f FilesBySequence) Less(i, j int) bool { return f[i].Metadata.Version < f[j].Metadata.Version }
func (f FilesBySequence) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func GetMigrationFiles(dir string) (*[]File, error) {
	migrationFiles := make([]File, 0)

	if !common.PathExist(dir) {
		return &migrationFiles, nil
	}

	err := filepath.Walk(dir,
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

			for _, file := range migrationFiles {
				if migrationFile.Metadata.Version == file.Metadata.Version {
					return fmt.Errorf("migration file '%s' and '%s' have the same version number (%d)",
						filepath.Base(migrationFile.Metadata.FileName),
						filepath.Base(file.Metadata.FileName),
						file.Metadata.Version)
				}
			}

			migrationFiles = append(migrationFiles, *migrationFile)

			return nil
		})
	if err != nil {
		return nil, err
	}

	sort.Sort(FilesBySequence(migrationFiles))

	return &migrationFiles, err
}

func parseMigrationFile(filename string, source []byte) (*File, error) {
	fileSha256 := fmt.Sprintf("%x", sha256.Sum256(source))

	sequenceNumber, err := getSequenceNumberFromFilename(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequence number from migration filename: %s, err: %s", filename, err)
	}

	var f FileContent
	err = hclsimple.Decode(filename, source, nil, &f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode migration file: %s, err: %s", filename, err)
	}

	var fileMeta = FileMetadata{
		FileName: filename,
		FileHash: fileSha256,
		Version:  sequenceNumber,
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

func getSequenceNumberFromFilename(filename string) (int, error) {
	regex := regexp.MustCompile(`V(?P<sequence>\d+)__`)
	match := regex.FindStringSubmatch(filename)

	if match == nil {
		return -1, fmt.Errorf("missing version prefix on migration file: '%s'", filename)
	}

	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return strconv.Atoi(result["sequence"])
}
