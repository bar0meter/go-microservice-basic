package configs

import (
	"os"
	"path/filepath"
)

// SendGridConfig => holds all the sendgrid required configurations
type SendGridConfig struct {
	APIKey string
}

// ServerConfig => Has all the servers configs (API keys, Client Secret, etc)
type ServerConfig struct {
	SendGrid *SendGridConfig
	RootPath string
	Providers *Providers
}

// Providers => Default notifications providers (Email,SMS) for server
type Providers struct {
	Email string
}

// NewConfig returns a new Config struct
// Has all the configs/credentials needed for all the services for this server
func NewConfig() *ServerConfig {
	sendGrid := NewSendGridConfig()
	rootPath, _ := filepath.Abs("./")
	providers := NewProviders()

	return &ServerConfig{
		SendGrid: sendGrid,
		RootPath: rootPath,
		Providers: providers,
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

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
