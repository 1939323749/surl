package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"surl/database"
	"time"
)

func ClickHandler(app *fiber.App) {
	app.Get("/click/:shortUrl", func(ctx *fiber.Ctx) error {
		shorturl := ctx.Params("shortUrl")
		collection := database.Db.Collection("urls")
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"shortUrl": shorturl}
		var result UrlMapping
		err := collection.FindOne(c, filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				log.Printf("Error finding short URL: %s", err)
				return err
			}
			ctx.Status(http.StatusInternalServerError)
			_, err = ctx.Writef("Error finding short URL: %s", err)
			if err != nil {
				return err
			}
		}
		ctx.Status(http.StatusOK)
		msg, _ := jsoniter.Marshal(UrlMapping{ShortUrl: result.ShortUrl, LongUrl: result.LongUrl, ClickCount: result.ClickCount})
		_, err = ctx.Write(msg)
		if err != nil {
			return err
		}
		return nil
	})
}
