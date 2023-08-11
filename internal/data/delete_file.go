package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) DeleteFile(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	var fileSize int
	fileSize, err = s.FileSize(ctx, id)

	if err != nil {
		return err
	}

	if fileSize != 0 { // uploaded
		err = s.bucket.DeleteContext(ctx, oid)

		if err != nil {
			return err
		}
	}

	_, err = s.allocatedFiles.DeleteOne(ctx, bson.M{"_id": oid})

	return err
}
