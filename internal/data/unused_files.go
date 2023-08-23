package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) UnusedFiles(ctx context.Context, until, limit int64) (ids []string, err error) {
	var cursor *mongo.Cursor
	opts := options.Find().SetProjection(bson.M{"_id": 1}).SetLimit(limit)
	cursor, err = s.unusedFiles.Find(ctx, bson.M{"since": bson.M{"$lte": until}}, opts)

	if err != nil {
		return nil, err
	}

	ids = make([]string, 0, limit)

	type Info struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	doc := &Info{}
	for cursor.Next(ctx) {
		err = cursor.Decode(doc)

		if err != nil {
			return nil, err
		}

		ids = append(ids, doc.Id.Hex())
	}

	return ids, nil
}
