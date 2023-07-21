package handler

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"surl/database"
	"time"
)

type body struct {
	Url string `json:"url"`
}

func CreateShortUrlHandler(app *fiber.App) {
	app.Post("/create", func(ctx *fiber.Ctx) error {
		reqBody := new(body)
		if err := ctx.BodyParser(reqBody); err != nil {
			return fmt.Errorf("invalid request body")
		}
		if !ValidUrl(reqBody.Url) {
			return fmt.Errorf("url can't be null")
		}
		collection := database.Db.Collection("urls")
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{"longUrl": reqBody.Url}
		var result UrlMapping
		err := collection.FindOne(c, filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				log.Printf("Error finding short URL: %s", err)
			} else {
				log.Printf("Error finding short URL: %s", err)
				ctx.Status(http.StatusInternalServerError)
				_, err = ctx.Writef("Error finding short URL: %s", err)
				if err != nil {
					return err
				}
				return err
			}
		}
		shortUrl := genShortUrl(6)
		if result.ShortUrl == "" {
			ctx.Status(http.StatusCreated)
			_, err = collection.InsertOne(c, UrlMapping{ShortUrl: shortUrl, LongUrl: reqBody.Url, ClickCount: 0})
			if err != nil {
				log.Printf("Error inserting short URL: %s", err)
				return err
			}
			status := database.RedisClient.Set(c, shortUrl, reqBody.Url, 24*time.Hour)
			if status.Err() != nil {
				log.Printf("Error inserting short URL: %s", err)
				return err
			}
			msg, err := jsoniter.Marshal(UrlMapping{ShortUrl: shortUrl, LongUrl: reqBody.Url, ClickCount: result.ClickCount})
			if err != nil {
				return err
			}
			_, err = ctx.Writef(string(msg))
			if err != nil {
				return err
			}
			return nil
		}

		msg, err := jsoniter.Marshal(UrlMapping{ShortUrl: result.ShortUrl, LongUrl: reqBody.Url, ClickCount: result.ClickCount})
		if err != nil {
			return err
		}
		ctx.Status(http.StatusOK)
		_, err = ctx.Writef(string(msg))
		if err != nil {
			return err
		}
		return nil
	})
}
