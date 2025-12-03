package config

import "os"

type Config struct {
	HTTPPort   string
	GRPCPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	JWTSecret  string
}

func LoadConfig() *Config {
	return &Config{
		HTTPPort:   getEnv("HTTP_PORT", "8081"),
		GRPCPort:   getEnv("GRPC_PORT", "9091"),
		DBHost:     getEnv("USER_DB_HOST", "localhost"),
		DBPort:     getEnv("USER_DB_PORT", "5432"),
		DBUser:     getEnv("USER_DB_USER", "chatuser"),
		DBPassword: getEnv("USER_DB_PASSWORD", "chatpass123"),
		DBName:     getEnv("USER_DB_NAME", "chatapp"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		JWTSecret:  getEnv("JWT_SECRET", "your-super-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
