package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gocs/pensive/internal/router"
	"github.com/gocs/pensive/tmpl"
)

func main() {
	// sets the session cookie store key
	session := getEnv("SESSION_KEY", "soopa-shiikurrets")
	// sets the redis localhost and port
	redisAddr := getEnv("REDIS_ADDR", "localhost:6380")
	// sets the redis password
	redisPassword := getEnv("REDIS_PASSWORD", "")
	// sets the file store assign address
	weedAddr := getEnv("SEAWEED_SERVER_ADDR", "http://seaweedfs:9333")
	// sets the file store upload address
	weedUpAddr := getEnv("SEAWEED_UPLOAD_ADDR", "http://seaweedfs:8080")
	// sets the file store upload address
	weedUpIP := getEnv("SEAWEED_UPLOAD_IP", "127.0.0.1")

	r, err := router.New(
		session,
		redisAddr,
		redisPassword,
		weedAddr,
		weedUpAddr,
		weedUpIP)
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
