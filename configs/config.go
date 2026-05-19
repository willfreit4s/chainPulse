package configs

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type conf struct {
	// Database configuration
	DBHost  string
	DBPort  int
	DBUser  string
	DBPass  string
	DBName  string
	MaxConn int
	MinConn int

	// Server configuration
	ServerPort     int
	ServiceName    string
	Environment    string
	ServiceVersion string

	// Auth configuration
	Auth0Domain   string
	Auth0Audience string
}

type Config = conf

func LoadConfig() (*conf, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	v := viper.New()

	// Defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASS", "postgres")
	v.SetDefault("DB_NAME", "chainpulse")
	v.SetDefault("MAX_CONN", 60)
	v.SetDefault("MIN_CONN", 20)
	v.SetDefault("SERVER_PORT", 8000)
	v.SetDefault("SERVICE_NAME", "chainpulse-api")
	v.SetDefault("SERVICE_VERSION", "1.0.0")
	v.SetDefault("ENVIRONMENT", "local")
	v.SetDefault("AUTH0_DOMAIN", "")
	v.SetDefault("AUTH0_AUDIENCE", "")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	cfg := &conf{
		DBHost:         v.GetString("DB_HOST"),
		DBPort:         v.GetInt("DB_PORT"),
		DBUser:         v.GetString("DB_USER"),
		DBPass:         v.GetString("DB_PASS"),
		DBName:         v.GetString("DB_NAME"),
		MaxConn:        v.GetInt("MAX_CONN"),
		MinConn:        v.GetInt("MIN_CONN"),
		ServerPort:     v.GetInt("SERVER_PORT"),
		ServiceName:    v.GetString("SERVICE_NAME"),
		ServiceVersion: v.GetString("SERVICE_VERSION"),
		Environment:    v.GetString("ENVIRONMENT"),
		Auth0Domain:    v.GetString("AUTH0_DOMAIN"),
		Auth0Audience:  v.GetString("AUTH0_AUDIENCE"),
	}

	return cfg, nil
}
