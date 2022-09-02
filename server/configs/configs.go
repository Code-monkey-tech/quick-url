package configs

import (
	"github.com/rs/zerolog/log"
	"os"
)

type AppConfig struct {
	RedisURI      string
	PostgresURL   string
	RedisPassword string
	ServerPort    string
}

func ReadConfig() *AppConfig {
	var conf = &AppConfig{
		PostgresURL:   mustEnv("PG_DB"),
		RedisURI:      mustEnv("RDB_URI"),
		RedisPassword: mustEnv("RDB_PASS"),
		ServerPort:    mustEnv("PORT"),
	}
	return conf
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Str("required parameter is not read from env", key)
	}
	return val
}
