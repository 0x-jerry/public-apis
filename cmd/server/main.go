package main

import (
	"log"
	"net/http"
	"os"

	"public-apis/internal/app"
	"public-apis/internal/config"
)

func main() {
	cfg := config.Load()
	cfg.PrintEnv()

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	handler := app.New(cfg)

	addr := host + ":" + port
	log.Printf("http://%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
