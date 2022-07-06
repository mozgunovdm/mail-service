package mongo

import (
	"context"
	"fmt"
	"mts/task-service/internal/config"
	"mts/task-service/internal/db"
	"time"

	"github.com/rs/zerolog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Collection *mongo.Collection
	Client     *mongo.Client
}

func NewMongo(ctx context.Context, cfg *config.Config, log *zerolog.Logger) (db.IDatabase, error) {

	log.Debug().Timestamp().Msg("Connecting to MongoDB")

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%v/%s",
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
		cfg.Mongo.Name)
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal().Err(err)
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err)
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err)
		return nil, err
	}

	collection := client.Database(cfg.Mongo.Name).Collection(cfg.Mongo.Collection)
	log.Debug().Timestamp().Msgf("MongoDB database: %v", collection.Name())

	m := &Mongo{
		Collection: collection,
		Client:     client,
	}

	return m, nil
}
