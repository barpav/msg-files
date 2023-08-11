package data

import (
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) DownloadFile(id string, stream io.Writer) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = s.bucket.DownloadToStream(oid, stream)

	return err
}
