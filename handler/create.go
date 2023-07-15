package handler

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"surl/database"
	"time"
)

// CreateShortUrlHandler handle /create
func CreateShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	longUrl := r.FormValue("url")
	if longUrl == "" {
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}
	collection := database.Db.Collection("urls")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingMapping UrlMapping
	err := collection.FindOne(ctx, bson.M{"longUrl": longUrl}).Decode(&existingMapping)
	if err == nil {
		jsonResponse, _ := json.Marshal(map[string]string{"shortUrl": existingMapping.ShortUrl})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		return
	}

	shortUrl := genShortUrl(6)

	_, err = collection.InsertOne(ctx, UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl, ClickCount: 0})
	if err != nil {
		log.Printf("Error inserting short URL: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	jsonResponse, _ := json.Marshal(map[string]string{"shortUrl": shortUrl})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
