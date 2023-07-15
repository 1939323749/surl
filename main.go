package main

import (
	"log"
	"net/http"
	"surl/database"
	"surl/handler"
)

func main() {
	database.Connect()
	http.HandleFunc("/", handler.RedirectHandler)
	http.HandleFunc("/create", handler.CreateShortUrlHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
