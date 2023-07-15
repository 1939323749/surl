package handler

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"surl/database"
	"time"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]
	if shortUrl == "create" {
		return
	}
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
	http.Redirect(w, r, result.LongUrl, http.StatusSeeOther)

	update := bson.M{
		"$inc": bson.M{"clickCount": 1},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error increasing click-count: %s", err)
	}
	//debug
	//err = collection.FindOne(ctx, filter).Decode(&result)
	//fmt.Println(result)
}
