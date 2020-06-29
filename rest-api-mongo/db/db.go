package db

import (
	"context"
	log "github.com/frost060/go-microservice-basic/rest-api-mongo/logging"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
	log.Info("Connecting to mongodb...")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Error("Error occurred while creating mongo client", "error", err)
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Error("Error occurred while connecting to mongodb", "error", err)
		return nil
	}

	log.Info("Successfully connected to mongodb...")
	return client.Database("go_mongo_todo")
}
