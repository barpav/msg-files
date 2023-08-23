package statistics

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FileUsage struct {
	FileId string `json:"fileId"`
	InUse  bool   `json:"inUse"`
}

func (f *FileUsage) serialize() ([]byte, error) {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(f)

	if err != nil {
		return nil, fmt.Errorf("failed to serialize file usage info: %w", err)
	}

	return b.Bytes(), nil
}

func (f *FileUsage) deserialize(data []byte) error {
	b := bytes.NewBuffer(data)
	err := json.NewDecoder(b).Decode(f)

	if err != nil {
		return fmt.Errorf("failed to deserialize file usage info: %w", err)
	}

	return nil
}
