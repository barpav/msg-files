package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) UpdateStats(ctx context.Context, fileId string, inUse bool) (uses int, err error) {
	var oid primitive.ObjectID
	oid, err = primitive.ObjectIDFromHex(fileId)

	if err != nil {
		return 0, err
	}

	update := 1
	if !inUse {
		update = -1
	}

	opts := options.FindOneAndUpdate().SetProjection(bson.M{"uses": 1})
	opts.SetReturnDocument(options.ReturnDocument(options.After))

	result := s.allocatedFiles.FindOneAndUpdate(ctx, bson.M{"_id": oid},
		bson.M{"$inc": bson.M{"uses": update}}, opts)

	info := &struct{ Uses int }{}
	err = result.Decode(info)

	return info.Uses, err
}
