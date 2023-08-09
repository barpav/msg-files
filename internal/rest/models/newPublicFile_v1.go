package models

import (
	"encoding/json"
	"errors"
	"io"
)

// Schema: newPublicFile.v1
type NewPublicFileV1 struct {
	Name string
	Mime string
}

func (m *NewPublicFileV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New file description violates 'newPublicFile.v1' schema.")
	}

	return m.validate()
}

func (m *NewPublicFileV1) validate() (err error) {
	if m.Name == "" {
		err = errors.Join(err, errors.New("File name must be specified."))
	}

	if m.Mime == "" {
		err = errors.Join(err, errors.New("MIME type must be specified."))
	}

	return err
}
