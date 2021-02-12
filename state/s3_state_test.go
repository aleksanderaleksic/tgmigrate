package state

import (
	"crypto/md5"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func generateRandomMd5() string {
	uuidString, err := uuid.GenerateUUID()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(uuidString)))
}

func TestParentObject(t *testing.T) {
	baseTime := time.Now()
	key := aws.String("file1/terraform.tfstate")

	latest := s3.ObjectVersion{
		IsLatest:     aws.Bool(true),
		Key:          key,
		LastModified: aws.Time(baseTime.Add(time.Duration(3 * 60 * 10000000))),
		VersionId:    aws.String(generateRandomMd5()),
	}
	expected := s3.ObjectVersion{
		IsLatest:     aws.Bool(false),
		Key:          key,
		LastModified: aws.Time(baseTime.Add(time.Duration(2 * 60 * 10000000))),
		VersionId:    aws.String(generateRandomMd5()),
	}
	testSlice := []*s3.ObjectVersion{
		&expected,
		{
			IsLatest:     aws.Bool(false),
			Key:          key,
			LastModified: aws.Time(baseTime.Add(time.Duration(1 * 60 * 10000000))),
			VersionId:    aws.String(generateRandomMd5()),
		},
		{
			IsLatest:     aws.Bool(false),
			Key:          key,
			LastModified: aws.Time(baseTime),
			VersionId:    aws.String(generateRandomMd5()),
		},
	}

	actual := parentObject(testSlice, latest)
	assert.Equal(t, expected, *actual)
}

func TestNoParentObject(t *testing.T) {
	baseTime := aws.Time(time.Now())
	key := aws.String("file3/terraform.tfstate")

	latest := s3.ObjectVersion{
		IsLatest:     aws.Bool(true),
		Key:          key,
		LastModified: baseTime,
		VersionId:    aws.String(generateRandomMd5()),
	}
	testSlice := []*s3.ObjectVersion{
		{
			IsLatest:     aws.Bool(false),
			Key:          aws.String("file1/terraform.tfstate"),
			LastModified: baseTime,
			VersionId:    aws.String(generateRandomMd5()),
		},
		{
			IsLatest:     aws.Bool(false),
			Key:          aws.String("file2/terraform.tfstate"),
			LastModified: baseTime,
			VersionId:    aws.String(generateRandomMd5()),
		},
	}

	actual := parentObject(testSlice, latest)
	assert.Nil(t, actual)
}

func TestChangedObjects(t *testing.T) {
	baseTime := aws.Time(time.Now())

	beforeSlice := []*s3.ObjectVersion{
		{
			IsLatest:     aws.Bool(true),
			Key:          aws.String("file1/terraform.tfstate"),
			LastModified: baseTime,
			VersionId:    aws.String(generateRandomMd5()),
		},
		{
			IsLatest:     aws.Bool(true),
			Key:          aws.String("file2/terraform.tfstate"),
			LastModified: baseTime,
			VersionId:    aws.String(generateRandomMd5()),
		},
	}

	afterSlice := []*s3.ObjectVersion{
		{
			IsLatest:     aws.Bool(true),
			Key:          aws.String("file1/terraform.tfstate"),
			LastModified: aws.Time(baseTime.Add(time.Duration(5 * 60 * 10000000))),
			VersionId:    aws.String(generateRandomMd5()),
		},
		{
			IsLatest:     aws.Bool(true),
			Key:          aws.String("file2/terraform.tfstate"),
			LastModified: aws.Time(baseTime.Add(time.Duration(5 * 60 * 10000000))),
			VersionId:    aws.String(generateRandomMd5()),
		},
		func() *s3.ObjectVersion {
			var obj = *beforeSlice[0]
			obj.IsLatest = aws.Bool(false)
			return &obj
		}(),
		func() *s3.ObjectVersion {
			var obj = *beforeSlice[1]
			obj.IsLatest = aws.Bool(false)
			return &obj
		}(),
		{
			IsLatest:     aws.Bool(true),
			Key:          aws.String("file3/terraform.tfstate"),
			LastModified: baseTime,
			VersionId:    aws.String(generateRandomMd5()),
		},
	}

	expected := []history.ChangedS3Object{
		{
			Key:           "file1/terraform.tfstate",
			FromVersionId: afterSlice[2].VersionId,
			ToVersionId:   *afterSlice[0].VersionId,
		},
		{
			Key:           "file2/terraform.tfstate",
			FromVersionId: afterSlice[3].VersionId,
			ToVersionId:   *afterSlice[1].VersionId,
		},
		{
			Key:           "file3/terraform.tfstate",
			FromVersionId: nil,
			ToVersionId:   *afterSlice[4].VersionId,
		},
	}

	actual := changedObjects(beforeSlice, afterSlice)

	assert.Equal(t, expected, actual)
}
