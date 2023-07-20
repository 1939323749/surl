package database

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	Db          *mongo.Database
	RedisClient *redis.Client
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
	err = retry.Do(func() error {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})
		_, err := RedisClient.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}, retry.Attempts(5), retry.Delay(2*time.Second))
	if err != nil {
		return err
	}
	return nil
}
