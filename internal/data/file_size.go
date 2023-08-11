package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) FileSize(ctx context.Context, id string) (size int, err error) {
	var oid primitive.ObjectID
	oid, err = primitive.ObjectIDFromHex(id)

	if err != nil {
		return 0, nil
	}

	type Result struct {
		Length int
	}

	var r Result
	opts := options.FindOne().SetProjection(bson.M{"length": 1})
	err = s.files.FindOne(ctx, bson.M{"_id": oid}, opts).Decode(&r)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}

		return 0, err
	}

	return r.Length, nil
}
