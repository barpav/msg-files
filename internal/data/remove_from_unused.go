package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) RemoveFromUnused(ctx context.Context, fileId string) error {
	oid, err := primitive.ObjectIDFromHex(fileId)

	if err == nil {
		_, err = s.unusedFiles.DeleteOne(ctx, bson.M{"_id": oid})
	}

	return err
}
