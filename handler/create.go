package handler

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
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

		shortUrl, err := database.RedisClient.Get(context.Background(), longUrl.Url).Result()
		if err == redis.Nil {
			shortUrl = genShortUrl(6)

			err = database.RedisClient.Set(context.Background(), longUrl.Url, shortUrl, 0).Err()
			if err != nil {
				return err
			}

			go func() {
				collection := database.Db.Collection("urls")
				c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err = collection.InsertOne(c, UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl.Url, ClickCount: 0})
				if err != nil {
					log.Printf("Error inserting short URL: %s", err)
				}
			}()
		} else if err != nil {
			return err
		}

		msg, err := jsoniter.Marshal(UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl.Url, ClickCount: 0})
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
