package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpireTimeInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://127.0.0.1"),
		Port:                   getEnv("PORT", "8081"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "admin"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "SHP"),
		JWTSecret:              getEnv("JWT_SECRET", "XAXAXAXA"),
		JWTExpireTimeInSeconds: getEnvAsInt("JWT_EXPIRE_TIME_IN_SECONDS", 3600*24*7),
	}
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
