package database

import (
	"testing"
	"time"

	"vrcomandaapi/internal/config"
)

func TestNewMongoManagerAndDBDefault(t *testing.T) {
	m := NewMongoManager(config.Config{MongoURI: "mongodb://127.0.0.1:1", MongoDatabase: "x"})
	if m == nil {
		t.Fatal("expected manager instance")
	}
	if m.DB() != nil {
		t.Fatal("expected nil db before connection")
	}
}

func TestInvalidateConnectionWhenNil(t *testing.T) {
	m := NewMongoManager(config.Config{})
	m.InvalidateConnection()
	if m.DB() != nil {
		t.Fatal("db should remain nil")
	}
}

func TestIsConnectedWhenNil(t *testing.T) {
	m := NewMongoManager(config.Config{})
	if m.isConnected() {
		t.Fatal("expected false when db is nil")
	}
}

func TestStartWithInvalidConfigKeepsNil(t *testing.T) {
	m := NewMongoManager(config.Config{MongoURI: "mongodb://127.0.0.1:1", MongoDatabase: "x"})
	m.Start(10 * time.Millisecond)

	time.Sleep(30 * time.Millisecond)
	if m.DB() != nil {
		t.Fatal("expected db to remain nil for unreachable mongo")
	}
}
