package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gocs/pensive/internal/router"
	"github.com/gocs/pensive/tmpl"
)

func main() {
	ctx := context.Background()

	r, err := router.New(ctx, &router.Config{
		// sets the session cookie store key
		SessionKey: getEnv("SESSION_KEY", "soopa-shiikurrets"),
		// sets the redis localhost and port
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6380"),
		// sets the redis password
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		// sets the gmail username as sender
		GmailEmail: getEnv("GMAIL_EMAIL", "example@example.com"),
		// sets the gmail app password of sender
		GmailPassword: getEnv("GMAIL_APP_PASSWORD", ""),
		// sets the jwt secret key
		AccessSecret: getEnv("ACCESS_SECRET", "soopa-shiikurrets-too"),
		// sets the minio api endpoint
		MinioEndpoint: getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
		// sets the minio username
		MinioUser: getEnv("MINIO_ROOT_USER", "minio"),
		// sets the minio password or API Key
		MinioPassword: getEnv("MINIO_ROOT_PASSWORD", "awaawawaaawawa123123xqcCursed"),
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
