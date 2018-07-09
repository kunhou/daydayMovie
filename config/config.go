package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Debug     bool
	Host      string
	JWTSecret string
	Level     string
	TMDBToken string
	DB        struct {
		Host      string
		Port      string
		DataBase  string
		User      string
		Password  string
		SSLEnable bool
		SSLPath   string
	}
	Google struct {
		ClientID     string
		ClientSecret string
	}
	Facebook struct {
		ClientID     string
		ClientSecret string
	}
}

var conf Config

func init() {
	conf.setDBsetting()
	conf.setLoginKey()
	conf.Debug = getEnvBool("DEBUG", false)
	conf.Host = getEnv("HOST")
	conf.JWTSecret = getEnv("JWT_SECRET")
	conf.Level = getEnv("LEVEL")
	conf.TMDBToken = getEnv("TMDB_TOKEN")
}

// GetConfig get all config
func GetConfig() *Config {
	return &conf
}

func (c *Config) setDBsetting() {
	c.DB.Host = getEnv("POSTGRES_HOST")
	c.DB.Port = getEnvWithDefault("POSTGRES_PORT", "5432")
	c.DB.DataBase = getEnv("POSTGRES_DBNAME")
	c.DB.User = getEnv("POSTGRES_USER")
	c.DB.Password = getEnv("POSTGRES_PASSWORD")
	c.DB.SSLEnable = getEnvBool("SSL_ENABLE", false)
	c.DB.SSLPath = getEnvWithDefault("POSTGRES_SSL", "")
}

func (c *Config) setLoginKey() {
	c.Google.ClientID = getEnv("GOOGLE_CLIENT_ID")
	c.Google.ClientSecret = getEnv("GOOGLE_CLIENT_SECRET")
	c.Facebook.ClientID = getEnv("FACEBOOK_CLIENT_ID")
	c.Facebook.ClientSecret = getEnv("FACEBOOK_CLIENT_SECRET")
}

func getEnvWithDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(key + " are not set")
}

func getEnvNumber(key string) int {
	if value, ok := os.LookupEnv(key); ok {
		c, err := strconv.Atoi(value)
		if err != nil {
			panic(key + " invalid")
		}
		return c
	}
	panic(key + " are not set")
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if strings.EqualFold(value, "true") {
			return true
		}
		return false
	}
	return fallback
}
