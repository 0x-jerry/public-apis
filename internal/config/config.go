package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AuthEnabled   bool
	AuthToken     string
	BrowserWS     string
	BrowserWSEnabled bool
	StoragePath   string
	OCRToken      string
}

func Load() *Config {
	return &Config{
		AuthEnabled:      isTrue(os.Getenv("AUTH_ENABLED")),
		AuthToken:        os.Getenv("AUTH_TOKEN"),
		BrowserWS:        getEnv("BROWSER_WS", "ws://127.0.0.1:9222"),
		BrowserWSEnabled: isTrue(os.Getenv("BROWSER_WS_ENABLED")),
		StoragePath:      getEnv("STORAGE_PATH", "files"),
		OCRToken:         os.Getenv("OCR_TOKEN"),
	}
}

func (c *Config) IsAuthConfigured() bool {
	return c.AuthEnabled && c.AuthToken != ""
}

func (c *Config) AbsStoragePath() string {
	abs, err := filepath.Abs(c.StoragePath)
	if err != nil {
		return c.StoragePath
	}
	return abs
}

func (c *Config) PrintEnv() {
	enabled := "false"
	if c.AuthEnabled {
		enabled = "true"
	}
	browserWS := "false"
	if c.BrowserWSEnabled {
		browserWS = "true"
	}
	println("Auth: ", enabled)
	println("Browser WS: ", browserWS)
	println("Storage: ", c.AbsStoragePath())
}

func isTrue(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "true" || v == "1"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
