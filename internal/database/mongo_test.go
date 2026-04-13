package database

import (
	"testing"
	"time"

	"vrcomandaapi/internal/config"
)

func TestConnectMongoInvalidURI(t *testing.T) {
	cfg := config.Config{
		MongoURI:      "://bad-uri",
		MongoDatabase: "vrcomanda",
	}

	if _, err := ConnectMongo(cfg); err == nil {
		t.Fatal("expected ConnectMongo to fail for invalid URI")
	}
}

func TestConnectMongoUnreachableServer(t *testing.T) {
	cfg := config.Config{
		MongoURI:      "mongodb://127.0.0.1:1",
		MongoDatabase: "vrcomanda",
	}

	start := time.Now()
	if _, err := ConnectMongo(cfg); err == nil {
		t.Fatal("expected ConnectMongo to fail for unreachable server")
	}
	if time.Since(start) > 12*time.Second {
		t.Fatal("ConnectMongo timeout took too long")
	}
}

func TestConnectMongoSuccess(t *testing.T) {
	cfg := config.Config{
		MongoURI:      "mongodb://127.0.0.1:27017",
		MongoDatabase: "vrcomanda_test",
	}

	db, err := ConnectMongo(cfg)
	if err != nil {
		t.Skipf("MongoDB não disponível, pulando teste de conexão bem-sucedida: %v", err)
	}

	if db == nil {
		t.Fatal("expected non-nil database handle on successful connection")
	}

	if db.Name() != cfg.MongoDatabase {
		t.Fatalf("expected database name %q, got %q", cfg.MongoDatabase, db.Name())
	}

	_ = db.Client().Disconnect(nil)
}
