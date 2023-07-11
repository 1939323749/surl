package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	_ "math/rand"
	"net/http"
	"time"
)

type UrlMapping struct {
	ShortUrl string `bson:"shortUrl"`
	LongUrl  string `bson:"longUrl"`
}

var db *mongo.Database

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("shorturl")

	http.HandleFunc("/", RedirectHandler)
	http.HandleFunc("/create", CreateShortUrlHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]
	if shortUrl == "create" {
		return
	}
	collection := db.Collection("urls")
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
	http.Redirect(w, r, result.LongUrl, http.StatusSeeOther)
}

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
	shortUrl := genShortUrl(6)

	collection := db.Collection("urls")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl})
	if err != nil {
		log.Printf("Error inserting short URL: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	jsonResponse, _ := json.Marshal(map[string]string{"shortUrl": shortUrl})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
func genShortUrl(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
