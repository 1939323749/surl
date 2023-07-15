package handler

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"strings"
	"surl/database"
	"time"
)

func ClickHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := strings.TrimPrefix(r.URL.Path, "/click/")
	//debug
	//fmt.Println(shortUrl)

	collection := database.Db.Collection("urls")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"shortUrl": shortUrl}
	var result UrlMapping
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
			return
		}
		log.Printf("Error finding short URL: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]string{"shortUrl": result.ShortUrl, "clickCount": strconv.Itoa(result.ClickCount)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
