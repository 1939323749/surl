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

// CreateShortUrlHandler handle /create
//
//	func CreateShortUrlHandler(w http.ResponseWriter, r *http.Request) {
//		if r.Method != "POST" {
//			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//			return
//		}
//		longUrl := r.FormValue("url")
//		if longUrl == "" {
//			http.Error(w, "Missing URL", http.StatusBadRequest)
//			return
//		}
//		collection := database.Db.Collection("urls")
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//
//		var existingMapping UrlMapping
//		err := collection.FindOne(ctx, bson.M{"longUrl": longUrl}).Decode(&existingMapping)
//		if err == nil {
//			jsonResponse, _ := json.Marshal(map[string]string{"shortUrl": existingMapping.ShortUrl})
//			w.Header().Set("Content-Type", "application/json")
//			w.Write(jsonResponse)
//			return
//		}
//
//		shortUrl := genShortUrl(6)
//
//		_, err = collection.InsertOne(ctx, UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl, ClickCount: 0})
//		if err != nil {
//			log.Printf("Error inserting short URL: %s", err)
//			http.Error(w, "Internal server error", http.StatusInternalServerError)
//			return
//		}
//		jsonResponse, _ := json.Marshal(map[string]string{"shortUrl": shortUrl})
//		w.Header().Set("Content-Type", "application/json")
//		w.Write(jsonResponse)
//	}
type shortUrl struct {
	Url string `json:"url"`
}

func CreateShortUrlHandler(app *fiber.App) {
	app.Post("/create", func(ctx *fiber.Ctx) error {
		longUrl := new(shortUrl)
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
			_, err := ctx.Writef("Internal server error", http.StatusInternalServerError)
			if err != nil {
				return err
			}
			return err
		}
		msg, err = jsoniter.Marshal(UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl.Url, ClickCount: 0})
		if err != nil {
			return err
		}
		_, err = ctx.Writef(string(msg))
		if err != nil {
			return err
		}
		return nil
	})
}
