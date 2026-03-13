package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config represents all application configuration
type Config struct {
	AppPort  int
	AppHost  string
	ENV      string
	Database DatabaseCfg
}

// DatabaseCfg contains database configuration
type DatabaseCfg struct {
	HOST string
	NAME string
	USER string
	PASS string
	PORT int
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
