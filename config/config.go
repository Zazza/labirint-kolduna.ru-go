package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config represents all application configuration
type Config struct {
	AppPort  int
	AppHost  string
	ENV      string
	Database DatabaseCfg
	JWT      JWTConfig
}

// DatabaseCfg contains database configuration
type DatabaseCfg struct {
	HOST string
	NAME string
	USER string
	PASS string
	PORT int
}

// JWTConfig contains JWT configuration
type JWTConfig struct {
	Secret        string
	Issuer        string
	AccessExpiry  string
	RefreshExpiry string
}

// NewConfig creates and returns validated configuration
func NewConfig() (*Config, error) {
	v := viper.New()
	//v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	// Check if environment variables are set
	if !hasEnvironmentVariables() {
		return nil, fmt.Errorf(
			"no environment variables found (APP_* prefix). " +
				"Make sure .env file was loaded or environment variables are set",
		)
	}

	// Build config
	cfg := buildConfig(v)

	return cfg, nil
}

// buildConfig constructs the Config struct from viper values
func buildConfig(v *viper.Viper) *Config {
	return &Config{
		AppPort: v.GetInt("APP_PORT"),
		AppHost: v.GetString("APP_HOST"),
		ENV:     v.GetString("APP_ENV"),
		Database: DatabaseCfg{
			HOST: v.GetString("DB_HOST"),
			NAME: v.GetString("DB_NAME"),
			USER: v.GetString("DB_USER"),
			PASS: v.GetString("DB_PASS"),
			PORT: v.GetInt("DB_PORT"),
		},
		JWT: JWTConfig{
			Secret:        v.GetString("JWT_SECRET"),
			Issuer:        v.GetString("JWT_ISSUER"),
			AccessExpiry:  v.GetString("JWT_ACCESS_EXPIRY"),
			RefreshExpiry: v.GetString("JWT_REFRESH_EXPIRY"),
		},
	}
}

// hasEnvironmentVariables checks if any APP_* variables are set
func hasEnvironmentVariables() bool {
	requiredVars := []string{
		"APP_HTTP_PORT",
		"APP_DATABASE_DSN",
		"APP_ENV",
	}

	for _, varName := range requiredVars {
		if os.Getenv(varName) != "" {
			return true
		}
	}

	return false
}

// LoadEnv loads environment variables from .env file
func LoadEnv(rootPath string) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	log.Printf("📋 Environment: %s\n", env)

	paths := []string{
		filepath.Join(rootPath, fmt.Sprintf(".env.%s", env)),
		filepath.Join(rootPath, ".env"),
	}

	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("✅ Loaded from: %s\n", path)
			return
		}
	}

	log.Printf("⚠️  No .env file found\n")
	log.Printf("ℹ️  Using system environment variables\n")
}

// Helper methods for environment checks

// IsTest returns true if environment is test
func (cfg *Config) IsTest() bool {
	return cfg.ENV == "test"
}

// IsDevelopment returns true if environment is development
func (cfg *Config) IsDevelopment() bool {
	return cfg.ENV == "dev"
}

// IsProduction returns true if environment is production
func (cfg *Config) IsProduction() bool {
	return cfg.ENV == "prod"
}

// ValidateJWTConfig validates JWT configuration
func ValidateJWTConfig(config *JWTConfig) error {
	if config.Secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if len(config.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters, got %d", len(config.Secret))
	}
	if config.Issuer == "" {
		return fmt.Errorf("JWT_ISSUER environment variable is required")
	}
	if config.AccessExpiry == "" {
		return fmt.Errorf("JWT_ACCESS_EXPIRY environment variable is required")
	}
	if _, err := time.ParseDuration(config.AccessExpiry); err != nil {
		return fmt.Errorf("JWT_ACCESS_EXPIRY must be a valid duration (e.g., '15m', '1h'): %w", err)
	}
	if config.RefreshExpiry == "" {
		return fmt.Errorf("JWT_REFRESH_EXPIRY environment variable is required")
	}
	if _, err := time.ParseDuration(config.RefreshExpiry); err != nil {
		return fmt.Errorf("JWT_REFRESH_EXPIRY must be a valid duration (e.g., '168h', '7d'): %w", err)
	}
	return nil
}

// ParseAccessExpiry parses and returns the access token expiry duration
func (cfg *JWTConfig) ParseAccessExpiry() (time.Duration, error) {
	if cfg.AccessExpiry == "" {
		return 15 * time.Minute, nil
	}
	return time.ParseDuration(cfg.AccessExpiry)
}

// ParseRefreshExpiry parses and returns the refresh token expiry duration
func (cfg *JWTConfig) ParseRefreshExpiry() (time.Duration, error) {
	if cfg.RefreshExpiry == "" {
		return 24 * 7 * time.Hour, nil
	}
	return time.ParseDuration(cfg.RefreshExpiry)
}
