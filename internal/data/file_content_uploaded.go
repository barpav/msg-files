package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) FileContentUploaded(ctx context.Context, id string) (bool, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	err = s.files.FindOne(ctx, bson.M{"_id": oid}).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
