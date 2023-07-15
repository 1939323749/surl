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
	http.HandleFunc("/click/", handler.ClickHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
