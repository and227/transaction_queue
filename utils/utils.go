package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func GetDbParams() (string, string) {
	dbDriver := "postgres"
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPortNum, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	dbDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPortNum, dbUser, dbPass, dbName,
	)

	return dbDriver, dbDSN
}

func GetRedisParams() (string, int, int) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
	}

	return redisHost, redisPort, redisDb
}
