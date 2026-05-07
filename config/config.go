package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
	Redis  RedisConfig
}

type DBConfig struct {
	DSNOverride string
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Addr string
	URL  string // для Render — полный URL
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExp, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	cfg := &Config{
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-me"),
			ExpirationHours: jwtExp,
		},
		Server: ServerConfig{
			// Render сам задаёт PORT
			Port: getEnv("PORT", getEnv("SERVER_PORT", "8080")),
		},
	}

	// Render подставляет DATABASE_URL автоматически
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DB = DBConfig{DSNOverride: dbURL}
	} else {
		cfg.DB = DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "workout_tracker"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		}
	}

	// Render подставляет REDIS_URL автоматически
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		cfg.Redis = RedisConfig{URL: redisURL}
	} else {
		cfg.Redis = RedisConfig{Addr: getEnv("REDIS_ADDR", "localhost:6379")}
	}

	return cfg, nil
}

func (c *DBConfig) DSN() string {
	if c.DSNOverride != "" {
		return c.DSNOverride
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
