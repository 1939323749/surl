package database

import (
	"context"
	"github.com/avast/retry-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	Db *mongo.Database
)

type Database interface {
	Connect()
}

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := retry.Do(func() error {
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
		if err != nil {
			log.Fatal(err)
			return nil
		}
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		Db = client.Database("shorturl")
		return nil
	}, retry.Attempts(5), retry.Delay(2*time.Second))
	if err != nil {
		return err
	}
	return nil
}
