package handler

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"surl/database"
	"time"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func RedirectHandler(app *fiber.App) {
	app.Get("/:shortUrl", func(ctx *fiber.Ctx) error {
		shortUrl := ctx.Params("shortUrl")
		fmt.Println(shortUrl)
		collection := database.Db.Collection("urls")
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"shortUrl": shortUrl}
		var result UrlMapping
		err := collection.FindOne(c, filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				_, err := ctx.Writef("%s", http.StatusBadRequest)
				if err != nil {
					return err
				}
			}
			return err
		}
		err = ctx.Redirect(result.LongUrl, http.StatusSeeOther)
		if err != nil {
			return err
		}

		update := bson.M{
			"$inc": bson.M{"clickCount": 1},
		}

		_, err = collection.UpdateOne(c, filter, update)
		if err != nil {
			log.Printf("Error increasing click-count: %s", err)
		}
		_, err = ctx.Write([]byte(shortUrl))
		if err != nil {
			return err
		}
		return nil
	})
}
