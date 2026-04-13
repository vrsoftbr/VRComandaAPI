package database

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"vrcomandaapi/internal/config"
)

// MongoManager holds the active Mongo database handle and manages reconnection in the background.
type MongoManager struct {
	cfg config.Config
	mu  sync.RWMutex
	db  *mongo.Database
}

var connectMongoFn = ConnectMongo

var pingMongoFn = func(ctx context.Context, db *mongo.Database) error {
	return db.Client().Ping(ctx, readpref.Primary())
}

var disconnectMongoFn = func(ctx context.Context, db *mongo.Database) error {
	return db.Client().Disconnect(ctx)
}

// NewMongoManager creates a new manager. Call Start to begin connection attempts.
func NewMongoManager(cfg config.Config) *MongoManager {
	return &MongoManager{cfg: cfg}
}

// DB returns the current Mongo database handle. Returns nil when not connected.
func (m *MongoManager) DB() *mongo.Database {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db
}

// InvalidateConnection invalidates the current cached handle after connection-level failures.
func (m *MongoManager) InvalidateConnection() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_ = disconnectMongoFn(ctx, m.db)
	cancel()

	m.db = nil
}

// Start attempts an immediate connection and launches a background goroutine
// that retries every interval while the database handle is unavailable
// or when the existing connection is no longer healthy.
func (m *MongoManager) Start(interval time.Duration) {
	m.connect()
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C
			if !m.isConnected() {
				m.connect()
			}
		}
	}()
}

func (m *MongoManager) isConnected() bool {
	m.mu.RLock()
	db := m.db
	m.mu.RUnlock()

	if db == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := pingMongoFn(ctx, db); err != nil {
		return false
	}

	return true
}

// connect opens a new MongoDB connection and stores the handle.
// Logs a warning on failure without affecting the running application.
func (m *MongoManager) connect() {
	db, err := connectMongoFn(m.cfg)
	if err != nil {
		slog.Warn("MongoDB indisponivel", "erro", err)
		return
	}

	m.mu.Lock()
	if m.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = disconnectMongoFn(ctx, m.db)
		cancel()
	}
	m.db = db
	m.mu.Unlock()
}

// IsMongoConnectionError checks whether an operation failed due to connection/topology issues.
func IsMongoConnectionError(err error) bool {
	if err == nil {
		return false
	}

	if mongo.IsNetworkError(err) {
		return true
	}

	var commandErr mongo.CommandError
	if errors.As(err, &commandErr) {
		if strings.EqualFold(commandErr.Name, "NotPrimaryOrSecondary") || strings.EqualFold(commandErr.Name, "NotPrimaryNoSecondaryOk") {
			return true
		}
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "server selection timeout") ||
		strings.Contains(message, "connection") ||
		strings.Contains(message, "topology") ||
		strings.Contains(message, "notprimaryorsecondary") ||
		strings.Contains(message, "not primary")
}
