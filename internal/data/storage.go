package data

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage struct {
	cfg    *Config
	client *mongo.Client
}

func (s *Storage) Open() (err error) {
	s.cfg = &Config{}
	s.cfg.Read()
	return s.connectToDatabase()
}

func (s *Storage) Close(ctx context.Context) (err error) {
	closed := make(chan struct{}, 1)

	go func() {
		err = s.client.Disconnect(ctx)
		closed <- struct{}{}
	}()

	select {
	case <-closed:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Storage) connectToDatabase() (err error) {
	dbAddress := fmt.Sprintf("mongodb://%s:%s", s.cfg.host, s.cfg.port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.client, err = mongo.Connect(ctx, options.Client().ApplyURI(dbAddress))

	if err == nil {
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err = s.client.Ping(ctx, readpref.Primary())
	}

	if err == nil {
		log.Info().Msg(fmt.Sprintf("Successfully connected to DB at %s", dbAddress))
	}

	return err
}
