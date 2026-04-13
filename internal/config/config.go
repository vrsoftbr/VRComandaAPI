package config

import "os"

// Config centralizes environment-driven runtime settings.
type Config struct {
	HTTPPort      string
	MongoURI      string
	MongoDatabase string
	SQLitePath    string
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		HTTPPort:      getEnv("HTTP_PORT"),
		MongoURI:      getEnv("MONGO_URI"),
		MongoDatabase: getEnv("MONGO_DATABASE"),
		SQLitePath:    getEnv("SQLITE_PATH"),
	}
}

// getEnv returns the raw environment value for a key.
func getEnv(key string) string {
	value := os.Getenv(key)

	return value
}
