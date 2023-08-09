package models

import (
	"encoding/json"
	"errors"
	"io"
)

// Schema: newPrivateFile.v1
type NewPrivateFileV1 struct {
	Name   string
	Mime   string
	Access []string
}

func (m *NewPrivateFileV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New file description violates 'newPrivateFile.v1' schema.")
	}

	return m.validate()
}

func (m *NewPrivateFileV1) validate() (err error) {
	if m.Name == "" {
		err = errors.Join(err, errors.New("File name must be specified."))
	}

	if m.Mime == "" {
		err = errors.Join(err, errors.New("MIME type must be specified."))
	}

	if len(m.Access) == 0 {
		err = errors.Join(err, errors.New("Users with access to the file must be specified."))
	}

	return err
}
