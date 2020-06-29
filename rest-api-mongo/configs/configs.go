package configs

import (
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleConfig struct {
	GoogleOauth *oauth2.Config
	RandomState string
}

type JWTConfig struct {
	SecretKey         []byte
	ExpirationTime    int64
	ExpireInThreshold int64
}

type SendGridConfig struct {
	ApiKey string
}

type Config struct {
	Google   *GoogleConfig
	JWT      *JWTConfig
	SendGrid *SendGridConfig
	RootPath string
}

// NewConfig returns a new Config struct
// Has all the configs/credentials needed for all the services for this server
func NewConfig() *Config {
	googleConfig := NewGoogleConfig()
	jwtConfig := NewJWTConfig()
	sendGrid := NewSendGridConfig()
	rootPath, _ := filepath.Abs("./")

	return &Config{
		Google:   googleConfig,
		JWT:      jwtConfig,
		SendGrid: sendGrid,
		RootPath: rootPath,
	}
}

func NewGoogleConfig() *GoogleConfig {
	googleOAuthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:9090/google/callback",
		ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	googleConfig := &GoogleConfig{googleOAuthConfig, "random"}
	return googleConfig
}

func NewJWTConfig() *JWTConfig {
	expTS := getEnv("JWT_EXPIRATION_TIME", "")
	expirationDuration, err := strconv.ParseInt(expTS, 10, 64)
	if expTS == "" || err != nil {
		expirationDuration = 5 * 60 // By default 5 mins expiration for JWT token
	}

	expInTS := getEnv("JWT_EXPIRE_IN_TIME", "")
	expireInDuration, err := strconv.ParseInt(expInTS, 10, 64)
	if expInTS == "" || err != nil {
		expirationDuration = 60 // By default refresh token only if it is about to expire in 1 min
	}

	return &JWTConfig{
		SecretKey:         []byte(getEnv("JWT_SECRET_KEY", "")),
		ExpirationTime:    expirationDuration,
		ExpireInThreshold: expireInDuration,
	}
}

func NewSendGridConfig() *SendGridConfig {
	apiKey := getEnv("SENDGRID_API_KEY", "")
	return &SendGridConfig{
		ApiKey: apiKey,
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
