package data

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AllocatedFile struct {
	Name   string
	Mime   string
	Access []string
}

func (s *Storage) AllocateNewFile(ctx context.Context, name, mime string, access []string) (id string, err error) {
	var result *mongo.InsertOneResult
	result, err = s.allocatedFiles.InsertOne(ctx, AllocatedFile{Name: name, Mime: mime, Access: access})

	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", errors.New("failed to receive inserted object id")
}
