package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

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
	m := NewMongoManager(config.Config{})

	connectCalls := 0
	connectMongoFn = func(_ config.Config) (*mongo.Database, error) {
		connectCalls++
		return nil, errors.New("unreachable")
	}
	pingMongoFn = func(_ context.Context, _ *mongo.Database) error {
		return errors.New("down")
	}

	m.Start(5 * time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	if m.DB() != nil {
		t.Fatal("expected db to remain nil when connect keeps failing")
	}
	if connectCalls < 2 {
		t.Fatalf("expected immediate and ticker connect attempts, got %d", connectCalls)
	}
}

func TestIsConnectedWithPingErrorAndSuccess(t *testing.T) {
	m := NewMongoManager(config.Config{})
	m.db = &mongo.Database{}

	originalPing := pingMongoFn
	t.Cleanup(func() {
		pingMongoFn = originalPing
	})

	pingMongoFn = func(_ context.Context, _ *mongo.Database) error {
		return errors.New("ping failed")
	}

	if m.isConnected() {
		t.Fatal("expected false when ping fails")
	}

	pingMongoFn = func(_ context.Context, _ *mongo.Database) error {
		return nil
	}

	if !m.isConnected() {
		t.Fatal("expected true when ping succeeds")
	}
}

func TestInvalidateConnectionWithActiveDB(t *testing.T) {
	m := NewMongoManager(config.Config{})
	m.db = &mongo.Database{}

	originalDisconnect := disconnectMongoFn
	t.Cleanup(func() {
		disconnectMongoFn = originalDisconnect
	})

	called := 0
	disconnectMongoFn = func(_ context.Context, _ *mongo.Database) error {
		called++
		return nil
	}

	m.InvalidateConnection()

	if called != 1 {
		t.Fatalf("expected disconnect to be called once, got %d", called)
	}
	if m.DB() != nil {
		t.Fatal("expected db to be cleared")
	}
}

func TestConnectHandlesErrorAndReplacesExistingDB(t *testing.T) {
	m := NewMongoManager(config.Config{})

	originalConnect := connectMongoFn
	originalDisconnect := disconnectMongoFn
	t.Cleanup(func() {
		connectMongoFn = originalConnect
		disconnectMongoFn = originalDisconnect
	})

	connectCalls := 0
	firstDB := &mongo.Database{}
	secondDB := &mongo.Database{}
	connectMongoFn = func(_ config.Config) (*mongo.Database, error) {
		connectCalls++
		if connectCalls == 1 {
			return nil, errors.New("connect error")
		}
		if connectCalls == 2 {
			return firstDB, nil
		}
		return secondDB, nil
	}

	disconnectCalls := 0
	disconnectMongoFn = func(_ context.Context, _ *mongo.Database) error {
		disconnectCalls++
		return nil
	}

	m.connect()
	if m.DB() != nil {
		t.Fatal("expected db to stay nil on connect error")
	}

	m.connect()
	if m.DB() != firstDB {
		t.Fatal("expected first successful connection to be stored")
	}

	m.connect()
	if m.DB() != secondDB {
		t.Fatal("expected db to be replaced on reconnect")
	}
	if disconnectCalls != 1 {
		t.Fatalf("expected previous db disconnect once, got %d", disconnectCalls)
	}
}

func TestIsMongoConnectionErrorAdditionalPaths(t *testing.T) {
	if !IsMongoConnectionError(mongo.CommandError{Labels: []string{"NetworkError"}}) {
		t.Fatal("expected labeled network error to be treated as connection error")
	}

	if !IsMongoConnectionError(mongo.CommandError{Name: "NotPrimaryNoSecondaryOk"}) {
		t.Fatal("expected NotPrimaryNoSecondaryOk command error to be treated as connection error")
	}

	if !IsMongoConnectionError(errors.New("topology changed")) {
		t.Fatal("expected topology message to be treated as connection error")
	}
}
