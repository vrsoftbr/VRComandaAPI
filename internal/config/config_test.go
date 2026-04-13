package config

import "testing"

func TestLoadReadsEnvironment(t *testing.T) {
	t.Setenv("HTTP_PORT", ":8080")
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")
	t.Setenv("MONGO_DATABASE", "vrcomanda")
	t.Setenv("SQLITE_PATH", "./data/app.db")

	cfg := Load()

	if cfg.HTTPPort != ":8080" {
		t.Fatalf("HTTPPort = %q", cfg.HTTPPort)
	}
	if cfg.MongoURI != "mongodb://localhost:27017" {
		t.Fatalf("MongoURI = %q", cfg.MongoURI)
	}
	if cfg.MongoDatabase != "vrcomanda" {
		t.Fatalf("MongoDatabase = %q", cfg.MongoDatabase)
	}
	if cfg.SQLitePath != "./data/app.db" {
		t.Fatalf("SQLitePath = %q", cfg.SQLitePath)
	}
}

func TestGetEnvReturnsEmptyWhenMissing(t *testing.T) {
	t.Setenv("MISSING_ENV_FOR_TEST", "")

	got := getEnv("MISSING_ENV_FOR_TEST")
	if got != "" {
		t.Fatalf("getEnv() = %q, want empty", got)
	}
}
