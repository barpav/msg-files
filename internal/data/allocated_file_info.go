package data

import (
	"context"

	"github.com/barpav/msg-files/internal/rest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) AllocatedFileInfo(ctx context.Context, id string) (info *models.AllocatedFile, err error) {
	var oid primitive.ObjectID
	oid, err = primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, nil
	}

	result := s.allocatedFiles.FindOne(ctx, bson.M{"_id": oid})
	err = result.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	info = &models.AllocatedFile{}
	err = result.Decode(info)

	if err != nil {
		return nil, err
	}

	return info, nil
}
