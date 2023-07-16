package handler

import "math/rand"

type UrlMapping struct {
	ShortUrl   string `bson:"shortUrl"`
	LongUrl    string `bson:"longUrl"`
	ClickCount int    `bson:"clickCount"`
}

func genShortUrl(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
func validUrl(url string) bool {
	//TODO
	return url != ""
}
