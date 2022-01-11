package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigType struct {
	RedisHost string
	RedisPort int
	Port      int
	SeedUrl   string

	DbName     string
	DbPassword string
	DbUser     string
	DbHost     string
	DbPort     int

	DequeueDelay    int
	CrawlerDuration int
}

var Config = &ConfigType{
	RedisHost:       "127.0.0.1",
	RedisPort:       6379,
	Port:            8000,
	SeedUrl:         "http://localhost:5000",
	DbName:          "webcrawler",
	DbPassword:      "",
	DbUser:          "",
	DbHost:          "0.0.0.0",
	DbPort:          3306,
	DequeueDelay:    50,
	CrawlerDuration: 1500,
}

// Load config from env to localConfig
func InitConfig() {

	if data := Getenv("REDIS_HOST"); data != "" {
		Config.RedisHost = data
	}
	if data := Getenv("REDIS_PORT"); data != "" {
		parsedData, err := strconv.Atoi(data)
		if err != nil {
			log.Printf("invalid data provided for REDIS_PORT, and got the error %v\n", err)
			return
		}
		Config.RedisPort = parsedData
	}

	if data := Getenv("PORT"); data != "" {
		parsedData, err := strconv.Atoi(data)
		if err != nil {
			log.Printf("Invalid data provided for PORT, and got the error %v\n", err)
			return
		}
		Config.Port = parsedData
	}
	if data := Getenv("SEED_URL"); data != "" {
		Config.SeedUrl = data
	}
	if data := Getenv("DB_NAME"); data != "" {
		Config.DbName = data
	}
	if data := Getenv("DB_PWD"); data != "" {
		Config.DbPassword = data
	}
	if data := Getenv("DB_USER"); data != "" {
		Config.DbUser = data
	}
	if data := Getenv("DB_HOST"); data != "" {
		Config.DbHost = data
	}

	if data := Getenv("DB_PORT"); data != "" {
		parsedData, err := strconv.Atoi(data)
		if err != nil {
			log.Printf("Invalid data provided for PORT, and got the error %v\n", err)
			return
		}
		Config.DbPort = parsedData
	}
	if data := Getenv("DEQUEUE_DELAY"); data != "" {
		parsedData, err := strconv.Atoi(data)
		if err != nil {
			log.Printf("Invalid data provided for PORT, and got the error %v\n", err)
			return
		}
		Config.DequeueDelay = parsedData
	}
	if data := Getenv("CRAWLER_DURATION"); data != "" {
		parsedData, err := strconv.Atoi(data)
		if err != nil {
			log.Printf("Invalid data provided for PORT, and got the error %v\n", err)
			return
		}
		Config.CrawlerDuration = parsedData
	}
}

// importing env
func Getenv(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("error loading .env file, %+v", err))
	}
	return os.Getenv(key)
}
