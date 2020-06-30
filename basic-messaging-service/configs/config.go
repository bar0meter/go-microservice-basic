package configs

import (
	"os"
	"path/filepath"
	"strconv"
)

// SendGridConfig => holds all the sendgrid required configurations
type SendGridConfig struct {
	APIKey string
}

// ServerConfig => Has all the servers configs (API keys, Client Secret, etc)
type ServerConfig struct {
	SendGrid  *SendGridConfig
	RootPath  string
	Providers *Providers
	Redis     *RedisConfig
}

// Providers => Default notifications providers (Email,SMS) for server
type Providers struct {
	Email string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// NewConfig returns a new Config struct
// Has all the configs/credentials needed for all the services for this server
func NewConfig() *ServerConfig {
	sendGrid := NewSendGridConfig()
	rootPath, _ := filepath.Abs("./")
	providers := NewProviders()
	redis := NewRedisConfig()

	return &ServerConfig{
		SendGrid:  sendGrid,
		RootPath:  rootPath,
		Providers: providers,
		Redis:     redis,
	}
}

func NewProviders() *Providers {
	email := getEnv("DEFAULT_EMAIL_PROVIDER", "sendgrid")
	return &Providers{
		Email: email,
	}
}

// NewSendGridConfig returns SendGrid configurations instance
func NewSendGridConfig() *SendGridConfig {
	apiKey := getEnv("SENDGRID_API_KEY", "")
	return &SendGridConfig{
		APIKey: apiKey,
	}
}

func NewRedisConfig() *RedisConfig {
	address := getEnv("REDIS_SERVER_ADDRESS", "localhost:6379")
	password := getEnv("REDIS_SERVER_PASSWORD", "")
	db, err := strconv.ParseInt(getEnv("REDIS_SERVER_DB", "0"), 10, 64)
	if err != nil {
		db = 0
	}

	return &RedisConfig{
		Addr:     address,
		Password: password,
		DB:       int(db),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
