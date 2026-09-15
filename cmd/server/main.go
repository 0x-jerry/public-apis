package main

import (
	"log"
	"net/http"
	"os"

	"public-apis/internal/app"
	"public-apis/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	cfg := config.Load()
	cfg.PrintEnv()

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	handler := app.New(cfg)

	addr := host + ":" + port
	log.Printf("http://%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

// loadEnv loads variables from a .env file when present. Existing environment
// variables take precedence and a missing file is ignored.
func loadEnv() {
	if _, err := os.Stat(".env"); err != nil {
		return
	}
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env: %v", err)
	}
}
