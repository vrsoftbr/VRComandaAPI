package database

import (
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

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

func TestIsMongoConnectionError(t *testing.T) {
	if IsMongoConnectionError(nil) {
		t.Fatal("nil error must return false")
	}

	if !IsMongoConnectionError(mongo.CommandError{Name: "NotPrimaryOrSecondary"}) {
		t.Fatal("expected command error to be treated as connection error")
	}

	if !IsMongoConnectionError(errors.New("server selection timeout")) {
		t.Fatal("expected timeout message to be treated as connection error")
	}

	if IsMongoConnectionError(errors.New("business validation failure")) {
		t.Fatal("non-connection error should return false")
	}
}
