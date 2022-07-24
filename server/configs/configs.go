package configs

import (
	"github.com/rs/zerolog/log"
	"os"
)

type AppConfig struct {
	RedisURI      string
	PostgresURL   string
	RedisPassword string
	ServerHost    string
	ServerPort    string
}

func ReadConfig() *AppConfig {
	var conf = &AppConfig{
		RedisURI:      mustEnv("rdb-uri"),
		PostgresURL:   mustEnv("pgdb-url"),
		RedisPassword: mustEnv("rdb-password"),
		ServerHost:    mustEnv("host"),
		ServerPort:    mustEnv("port"),
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
