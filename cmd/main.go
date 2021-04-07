package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gocs/pensive/pkg/router"
)

func main() {
	// sets the session cookie store key
	session := getEnv("SESSION_KEY", "soopa-shiikurrets")
	// sets the redis localhost and port
	redisAddr := getEnv("REDIS_ADDR", "localhost:6380")
	// sets the redis password
	redisPassword := getEnv("REDIS_PASSWORD", "")

	r, err := router.New(session, redisAddr, redisPassword)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
