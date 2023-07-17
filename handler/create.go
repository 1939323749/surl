package handler

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"surl/database"
	"time"
)

type longUrl struct {
	Url string `json:"url"`
}

func CreateShortUrlHandler(app *fiber.App) {
	app.Post("/create", func(ctx *fiber.Ctx) error {
		longUrl := new(longUrl)
		if err := ctx.BodyParser(longUrl); err != nil {
			return err
		}
		if !validUrl(longUrl.Url) {
			return fmt.Errorf("url can't be null")
		}
		collection := database.Db.Collection("urls")
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existingMapping UrlMapping
		err := collection.FindOne(c, bson.M{"longUrl": longUrl.Url}).Decode(&existingMapping)
		var msg []byte
		if err == nil {
			msg, err = jsoniter.Marshal(UrlMapping{ShortUrl: existingMapping.ShortUrl, LongUrl: existingMapping.LongUrl, ClickCount: existingMapping.ClickCount})
			ctx.Status(http.StatusOK)
			_, err := ctx.Writef(string(msg))
			if err != nil {
				return err
			}
			return nil
		}

		shortUrl := genShortUrl(6)

		_, err = collection.InsertOne(c, UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl.Url, ClickCount: 0})
		if err != nil {
			log.Printf("Error inserting short URL: %s", err)
			ctx.Status(http.StatusInternalServerError)
			_, err := ctx.Writef("%s", err)
			if err != nil {
				return err
			}
			return err
		}
		msg, err = jsoniter.Marshal(UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl.Url, ClickCount: 0})
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
