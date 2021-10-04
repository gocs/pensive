package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gocs/pensive/internal/router"
	"github.com/gocs/pensive/tmpl"
)

func main() {
	r, err := router.New(&router.Config{
		// sets the session cookie store key
		SessionKey: getEnv("SESSION_KEY", "soopa-shiikurrets"),
		// sets the redis localhost and port
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6380"),
		// sets the redis password
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		// sets the file store assign address
		WeedAddr: getEnv("SEAWEED_SERVER_ADDR", "http://seaweedfs:9333"),
		// sets the file store upload address
		WeedUpAddr: getEnv("SEAWEED_UPLOAD_ADDR", "http://seaweedfs:8080"),
		// sets the file store upload ip
		WeedUpIP: getEnv("SEAWEED_UPLOAD_IP", "http://127.0.0.1:8080"),
		// sets the gmail username as sender
		GmailEmail: getEnv("GMAIL_EMAIL", "example@example.com"),
		// sets the gmail app password of sender
		GmailPassword: getEnv("GMAIL_APP_PASSWORD", ""),
		// sets the jwt secret key
		AccessSecret: getEnv("ACCESS_SECRET", "soopa-shiikurrets-too"),
	})
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", tmpl.AssetsFS())
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
