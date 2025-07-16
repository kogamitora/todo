package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	DSN      string `json:"dsn"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", ""),
			Port: getEnv("SERVER_PORT", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", ""),
			DSN:      getEnv("DB_DSN", ""),
		},
	}

	return config, nil
}

// GetDSN builds the database connection string.
func (c *Config) GetDSN() string {
	if c.Database.DSN != "" {
		return c.Database.DSN
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	)
}

// GetServerAddr returns the full server address (host:port).
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("SERVER_HOST is required")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}
	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		return fmt.Errorf("invalid SERVER_PORT: %s", c.Server.Port)
	}

	// If DSN is not provided, validate individual DB components.
	if c.Database.DSN == "" {
		if c.Database.Host == "" {
			return fmt.Errorf("DB_HOST is required")
		}
		if c.Database.Port == "" {
			return fmt.Errorf("DB_PORT is required")
		}
		if c.Database.User == "" {
			return fmt.Errorf("DB_USER is required")
		}
		if c.Database.Password == "" {
			return fmt.Errorf("DB_PASSWORD is required")
		}
		if c.Database.Database == "" {
			return fmt.Errorf("DB_NAME is required")
		}
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadFromFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read env file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if they exist
		if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}

	return nil
}
