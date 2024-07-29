package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWTExpireTimeInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
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
