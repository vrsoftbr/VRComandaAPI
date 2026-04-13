package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"vrcomandaapi/internal/config"
)

func TestSwaggerHostFromPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "localhost:8080"},
		{name: "spaced", in: "  ", want: "localhost:8080"},
		{name: "port only", in: ":9090", want: "localhost:9090"},
		{name: "host and port", in: "127.0.0.1:8081", want: "127.0.0.1:8081"},
		{name: "numeric", in: "8082", want: "localhost:8082"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := swaggerHostFromPort(tc.in)
			if got != tc.want {
				t.Fatalf("swaggerHostFromPort(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestBootstrapSuccess(t *testing.T) {
	t.Setenv("HTTP_PORT", ":18080")
	t.Setenv("MONGO_URI", "mongodb://127.0.0.1:1")
	t.Setenv("MONGO_DATABASE", "vrcomanda")
	t.Setenv("SQLITE_PATH", ":memory:")

	router, err := bootstrap()
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestBootstrapSQLiteError(t *testing.T) {
	t.Setenv("HTTP_PORT", ":18080")
	t.Setenv("MONGO_URI", "mongodb://127.0.0.1:1")
	t.Setenv("MONGO_DATABASE", "vrcomanda")

	file := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to create file fixture: %v", err)
	}
	t.Setenv("SQLITE_PATH", filepath.Join(file, "db.sqlite"))

	_, err := bootstrap()
	if err == nil {
		t.Fatal("expected bootstrap to fail for invalid sqlite path")
	}
}

func TestBootstrapAutoMigrateError(t *testing.T) {
	originalAutoMigrate := autoMigrateSQLite
	defer func() { autoMigrateSQLite = originalAutoMigrate }()

	autoMigrateSQLite = func(_ *gorm.DB) error {
		return errors.New("migrate failure")
	}

	t.Setenv("HTTP_PORT", ":18080")
	t.Setenv("MONGO_URI", "mongodb://127.0.0.1:1")
	t.Setenv("MONGO_DATABASE", "vrcomanda")
	t.Setenv("SQLITE_PATH", ":memory:")

	_, err := bootstrap()
	if err == nil {
		t.Fatal("expected bootstrap to fail on auto-migrate error")
	}
}

func TestRunBootstrapError(t *testing.T) {
	errExpected := errors.New("bootstrap failure")
	err := run(
		func() (*gin.Engine, error) { return nil, errExpected },
		func() config.Config { return config.Config{HTTPPort: ":8080"} },
	)

	if !errors.Is(err, errExpected) {
		t.Fatalf("run error = %v, want %v", err, errExpected)
	}
}

func TestRunServerError(t *testing.T) {
	originalRunServer := runServer
	defer func() { runServer = originalRunServer }()

	runServer = func(_ *gin.Engine, _ string) error {
		return errors.New("listen failure")
	}

	r := gin.New()
	err := run(
		func() (*gin.Engine, error) { return r, nil },
		func() config.Config { return config.Config{HTTPPort: ":8080"} },
	)

	if err == nil {
		t.Fatal("expected run to fail when runServer returns error")
	}
}

func TestRunSuccessReturnsNil(t *testing.T) {
	originalRunServer := runServer
	defer func() { runServer = originalRunServer }()

	runServer = func(_ *gin.Engine, _ string) error {
		return nil
	}

	err := run(
		func() (*gin.Engine, error) { return gin.New(), nil },
		func() config.Config { return config.Config{HTTPPort: ":8080"} },
	)

	if err != nil {
		t.Fatalf("expected nil error on successful run, got: %v", err)
	}
}

func TestMainUsesFatalLogOnRunError(t *testing.T) {
	originalFatal := fatalLog
	defer func() { fatalLog = originalFatal }()

	called := false
	fatalLog = func(v ...any) { called = true }

	// Force bootstrap error by using a file as parent directory for SQLite path.
	file := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to create file fixture: %v", err)
	}
	t.Setenv("SQLITE_PATH", filepath.Join(file, "db.sqlite"))
	t.Setenv("HTTP_PORT", ":18080")
	t.Setenv("MONGO_URI", "mongodb://127.0.0.1:1")
	t.Setenv("MONGO_DATABASE", "vrcomanda")

	main()
	if !called {
		t.Fatal("expected fatalLog to be called")
	}
}
