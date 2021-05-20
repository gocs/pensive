package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gocs/pensive/internal/router"
)

//go:embed assets
var assets embed.FS

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

	r, err := router.New(session, redisAddr, redisPassword, weedAddr, weedUpAddr)
	if err != nil {
		log.Fatal(err)
	}

	// css and js files
	fs := http.FileServer(http.FS(assets))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
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
