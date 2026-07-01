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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	handler := app.New(cfg)

	log.Printf("http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
