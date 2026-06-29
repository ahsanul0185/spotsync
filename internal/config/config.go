package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

func LoadEnv() *Config {

	// .env file is optional — on production (e.g. Render), env vars are injected directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables from system")
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		Dsn:       os.Getenv("DSN"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
