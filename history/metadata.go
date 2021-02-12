package history

import (
	"encoding/json"
	"fmt"
)

type MetadataWrapper struct {
	S3Metadata *StorageS3Metadata
}

func (a *MetadataWrapper) MarshalJSON() ([]byte, error) {
	if a.S3Metadata != nil {
		return json.Marshal(a.S3Metadata)
	}
	return nil, fmt.Errorf("metadata is required")
}

func (a *MetadataWrapper) UnmarshalJSON(value []byte) error {
	var kv map[string]interface{}
	if err := json.Unmarshal(value, &kv); err != nil {
		return err
	}
	t := kv["type"].(string)
	switch t {
	case "s3":
		var meta StorageS3Metadata
		err := json.Unmarshal(value, &meta)
		a.S3Metadata = &meta
		return err
	default:
		return fmt.Errorf("unable to unmarshall MetadataWrapper for type: %s", t)
	}
}
