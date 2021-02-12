package state

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"sort"
)

type S3State struct {
	context   common.Context
	State     config.S3StateConfig
	Sync      S3Sync
	S3        *s3.S3
	Terraform *tfexec.Terraform
	Cache     common.Cache
}

func (s *S3State) InitializeState() error {
	err := os.MkdirAll(getStateDirPath(s.Cache), os.ModePerm)
	if err != nil {
		return err
	}

	tf, err := initializeTerraformExec(getStateDirPath(s.Cache))
	s.Terraform = tf
	if err != nil {
		return err
	}

	err = s.Sync.DownSync3State()
	if err != nil {
		return err
	}

	return nil
}

func (s S3State) Complete() (*history.MetadataWrapper, error) {
	beforeObjectVersions, err := s.ListAllObjects()
	if err != nil {
		return nil, err
	}

	err = s.Sync.UpSync3State()
	if err != nil {
		return nil, err
	}

	afterObjectVersions, err := s.ListAllObjects()
	if err != nil {
		return nil, err
	}

	changedObjects := changedObjects(beforeObjectVersions, afterObjectVersions)

	return history.StorageS3Metadata{
		SchemaVersion:  "v1",
		Type:           "s3",
		ChangedObjects: changedObjects,
	}.Wrap(), nil
}

func (s S3State) ListAllObjects() ([]*s3.ObjectVersion, error) {
	versions := make([]*s3.ObjectVersion, 0)

	err := s.S3.ListObjectVersionsPages(&s3.ListObjectVersionsInput{
		Bucket: aws.String(s.State.Bucket),
		Prefix: s.State.Prefix,
	}, func(output *s3.ListObjectVersionsOutput, b bool) bool {
		versions = append(versions, output.Versions...)
		return output.NextKeyMarker != nil //Continue until no more pages
	})
	return versions, err
}

func changedObjects(beforeSlice []*s3.ObjectVersion, afterSlice []*s3.ObjectVersion) []history.ChangedS3Object {
	newObjects := make([]*s3.ObjectVersion, 0)

	for _, afterObject := range afterSlice {
		found := false
		for _, beforeObject := range beforeSlice {
			if *afterObject.Key == *beforeObject.Key &&
				*afterObject.VersionId == *beforeObject.VersionId {
				found = true
				break
			}
		}

		if !found {
			newObjects = append(newObjects, afterObject)
		}
	}

	changedObjects := make([]history.ChangedS3Object, 0)
	for _, newObject := range newObjects {
		parent := parentObject(afterSlice, *newObject)
		if parent != nil {
			changedObjects = append(changedObjects, history.ChangedS3Object{
				Key:           *newObject.Key,
				FromVersionId: parent.VersionId,
				ToVersionId:   *newObject.VersionId,
			})
		} else {
			changedObjects = append(changedObjects, history.ChangedS3Object{
				Key:           *newObject.Key,
				FromVersionId: nil,
				ToVersionId:   *newObject.VersionId,
			})
		}

	}

	return changedObjects
}

func parentObject(slice []*s3.ObjectVersion, changedObject s3.ObjectVersion) *s3.ObjectVersion {
	var objectWithKey []*s3.ObjectVersion
	for _, obj := range slice {
		if *obj.Key == *changedObject.Key && *obj.VersionId != *changedObject.VersionId {
			objectWithKey = append(objectWithKey, obj)
		}
	}

	if len(objectWithKey) == 0 {
		return nil
	}

	sort.Slice(objectWithKey, func(i, j int) bool {
		return slice[i].LastModified.After(*slice[j].LastModified)
	})

	return objectWithKey[0]
}

func (s S3State) Move(from ResourceContext, to ResourceContext) (bool, error) {
	return move(s.Terraform, getStateDirPath(s.Cache), s.State.GetStateFileName(), from, to)
}

func (s S3State) Remove(resource ResourceContext) (bool, error) {
	return remove(s.Terraform, getStateDirPath(s.Cache), s.State.GetStateFileName(), resource)
}

func (s S3State) Cleanup() {
	os.RemoveAll(getStateDirPath(s.Cache))
}
