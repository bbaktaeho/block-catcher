package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_SCHEMA   string
}

var Env *Environment

func getEnv(key string, required bool) string {
	value, ok := os.LookupEnv(key)
	if !ok && required {
		log.Fatalf("Missing or invalid environment key: '%s'", key)
	}

	return value
}

func loadEnvironment() {
	if Env == nil {
		Env = new(Environment)
	}

	Env.DB_USER = getEnv("DB_USER", true)
	Env.DB_PASSWORD = getEnv("DB_PASSWORD", true)
	Env.DB_HOST = getEnv("DB_HOST", true)
	Env.DB_PORT = getEnv("DB_PORT", true)
	Env.DB_SCHEMA = getEnv("DB_SCHEMA", true)
}

func LoadEnvironmentFile(file string) {
	if err := godotenv.Load(file); err != nil {
		log.Fatalf("Error on load environment file: %s", file)
	}

	loadEnvironment()
}
