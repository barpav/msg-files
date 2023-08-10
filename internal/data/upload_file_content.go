package data

import (
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) UploadFileContent(id string, content io.Reader) (err error) {
	var oid primitive.ObjectID
	oid, err = primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	return s.bucket.UploadFromStreamWithID(oid, id, content)
}
