package migration

import (
	"github.com/aleksanderaleksic/tgmigrate/test"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGetMigrationFilesFromEmptyDirectory(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()

	files, err := GetMigrationFiles(testDir)
	ass.Nil(err)
	ass.Len(*files, 0)
}

func TestHandleInvalidMigrationFilesError(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	test.CopyTestData(t, "empty_migration", testDir)

	files, err := GetMigrationFiles(filepath.Join(testDir, "migrations"))
	ass.NotNil(err)
	ass.Nil(files)
}

func TestHandleInvalidMoveBlock(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	test.CopyTestData(t, "invalid_move_block", testDir)

	files, err := GetMigrationFiles(filepath.Join(testDir, "migrations"))
	ass.NotNil(err)
	ass.Nil(files)
}

func TestHandleInvalidRemoveBlock(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	test.CopyTestData(t, "invalid_remove_block", testDir)

	files, err := GetMigrationFiles(filepath.Join(testDir, "migrations"))
	ass.NotNil(err)
	ass.Nil(files)
}

func TestHandleMissingVersionPrefix(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	test.CopyTestData(t, "missing_version", testDir)

	files, err := GetMigrationFiles(filepath.Join(testDir, "migrations"))
	ass.NotNil(err)
	ass.Nil(files)
}

func TestHandleDuplicateVersionPrefix(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	test.CopyTestData(t, "duplicate_version", testDir)

	files, err := GetMigrationFiles(filepath.Join(testDir, "migrations"))
	ass.NotNil(err)
	ass.Nil(files)
}

func TestHandleFilePermissionDenied(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	test.CopyTestData(t, "duplicate_version", testDir)

	files, err := GetMigrationFiles(filepath.Join(testDir, "migrations"))
	ass.NotNil(err)
	ass.Nil(files)
}

func TestValidMigration(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	test.CopyTestData(t, "simple", testDir)

	f, err := GetMigrationFiles(filepath.Join(testDir, "migrations"))
	ass.Nil(err)

	files := *f
	ass.Len(files, 2)
	ass.Equal(files[0], File{
		Metadata: FileMetadata{
			FileName: "V1__move.hcl",
			FileHash: "5f18eb04c47a8f7ead44b7e727e1d21f9c07954c3143180383a6cc2280e2fc0d",
			Version:  1,
		},
		Config: Config{
			Environments: []string{"test"},
			Description:  "    - Move rest_api lambda to rest_v2\n",
		},
		Migrations: []Migration{
			{
				Type: "move",
				Name: "rest_api",
				Move: &MoveBlock{
					From: MoveFromBlock{
						State:    "us-east-1/apis/rest",
						Resource: "aws_lambda_function.rest_api",
					},
					To: MoveFromBlock{
						State:    "us-east-1/apis/rest_v2",
						Resource: "aws_lambda_function.rest_api",
					},
				},
				Remove: nil,
			},
		},
	})

	ass.Equal(files[1], File{
		Metadata: FileMetadata{
			FileName: "V2__remove.hcl",
			FileHash: "92cb5ae4c96519f4213a6ee2d0ada9c360866ea81cd7977f2caabbebd736b598",
			Version:  2,
		},
		Config: Config{
			Environments: []string{"test"},
			Description:  "    - Remove testfile from files module\n",
		},
		Migrations: []Migration{
			{
				Type: "remove",
				Name: "file",
				Move: nil,
				Remove: &RemoveBlock{
					State:    "us-east-1/files",
					Resource: "file.test_file",
				},
			},
		},
	})
}
