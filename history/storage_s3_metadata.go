package history

type StorageS3Metadata struct {
	SchemaVersion  string            `json:"schema_version"`
	Type           string            `json:"type"`
	ChangedObjects []ChangedS3Object `json:"changed_objects"`
}

type ChangedS3Object struct {
	Key           string  `json:"key"`
	FromVersionId *string `json:"from_version_id,optional"`
	ToVersionId   string  `json:"to_version_id"`
}

func (s StorageS3Metadata) Wrap() *MetadataWrapper {
	return &MetadataWrapper{
		S3Metadata: &s,
	}
}
