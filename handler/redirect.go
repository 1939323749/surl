package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"surl/database"
	"time"
)

func RedirectHandler(app *fiber.App) {
	app.Get("/:shortUrl", func(ctx *fiber.Ctx) error {
		shortUrl := ctx.Params("shortUrl")
		log.Printf("shortUrl: %s", shortUrl)
		collection := database.Db.Collection("urls")
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"shortUrl": shortUrl}
		var result UrlMapping
		err := collection.FindOne(c, filter).Decode(&result)
		if result.LongUrl == "" {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			err = ctx.Redirect(result.LongUrl, http.StatusSeeOther)
			if err != nil {
				log.Printf("Error redirecting: %s", err)
			}
		}
		if err != nil {
			log.Printf("Error finding short URL: %s", err)
		}
		updateCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		update := bson.M{
			"$inc": bson.M{"clickCount": 1},
		}
		_, err = collection.UpdateOne(updateCtx, filter, update)
		if err != nil {
			log.Printf("Error increasing click-count: %s", err)
		}
		return nil
	})
}
