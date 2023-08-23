package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) MarkAsUnused(ctx context.Context, fileId string) error {
	oid, err := primitive.ObjectIDFromHex(fileId)

	if err == nil {
		_, err = s.unusedFiles.InsertOne(ctx, bson.M{"_id": oid, "since": time.Now().UTC().Unix()})
	}

	return err
}
