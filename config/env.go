package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost         string
	Port               string
	DBUser             string
	DBPassword         string
	DBHost             string
	DBName             string
	JWTExpireInSeconds int64
	JWTSecret          string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:         getEnv("PUBLIC_HOST", "localhost"),
		Port:               getEnv("PORT", "8080"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", "root"),
		DBHost:             fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3306")),
		DBName:             getEnv("DB_NAME", "ecom"),
		JWTExpireInSeconds: getEnvsAsInt("JWT_EXPIRE_IN_SECONDS", 3600*24*7),
		JWTSecret:          getEnv("JWT_SECRET", "no-need-secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvsAsInt(key string, Fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Fallback
		}
		return i
	}

	return Fallback
}
