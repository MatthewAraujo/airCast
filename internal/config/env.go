package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type JWT struct {
	JWTExpirationInSeconds int64
	JWTSecret              string
}

type API struct {
	PublicHost string
	Port       string
}

type Postgres struct {
	URL      string
	Username string
	Password string
	Host     string
	Port     int64
	DBName   string
	SSLMode  string
}

type Config struct {
	JWT      JWT
	API      API
	Postgres Postgres
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		API: API{
			PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
			Port:       getEnv("PORT", "8080"),
		},
		JWT: JWT{
			JWTSecret:              getEnv("JWT_SECRET", "not-that-secret"),
			JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		},
		Postgres: Postgres{
			URL:      getEnv("POSTGRES_URL", "postgresql://docker:docker@pg:5432/airCast"),
			Username: getEnv("POSTGRES_USER", "docker"),
			Password: getEnv("POSTGRES_PASSWORD", "docker"),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			DBName:   getEnv("POSTGRES_DB", "airCast"),
			SSLMode:  getEnv("POSTGRES_SSLMode", "disabled"),
		},
	}

}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
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
