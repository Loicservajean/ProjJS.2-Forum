package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	errLoad := godotenv.Load("./.env")
	if errLoad != nil {
		log.Println("aucun fichier .env trouve : ")
	}
}

func GetEnvWithDefault(key string, defaultValue string) string {
	envVar, envErr := os.LookupEnv(key)
	if !envErr {
		return defaultValue
	}
	return envVar
}

func GetRequiredEnv(key string) string {
	envVar, envErr := os.LookupEnv(key)
	if !envErr {
		log.Fatalf("variable denvironnement manquante : %s", key)
	}
	return envVar
}
