package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Get() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Tidak terdapat menemukan file .env")
	}

	expInt, err := strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		log.Println("Kesalahan JWT_EXP")
	}

	redisDBIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDBIndex = 0
	}

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		MysqlDB: MysqlDB{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Name: os.Getenv("DB_NAME"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Tz:   os.Getenv("DB_TZ"),
		},
		RedisDB: RedisDB{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
			Pass: os.Getenv("REDIS_PASS"),
			DB:   redisDBIndex,
		},
		Jwt: Jwt{
			Key: os.Getenv("JWT_KEY"),
			Exp: expInt,
		},
	}
}
