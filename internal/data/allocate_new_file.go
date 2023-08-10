package data

import (
	"context"
	"errors"

	"github.com/barpav/msg-files/internal/rest/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) AllocateNewFile(ctx context.Context, fileInfo *models.AllocatedFile) (id string, err error) {
	var result *mongo.InsertOneResult
	result, err = s.allocatedFiles.InsertOne(ctx, fileInfo)

	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", errors.New("failed to receive inserted object id")
}
