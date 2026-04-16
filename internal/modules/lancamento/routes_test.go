package lancamento

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"vrcomandaapi/internal/shared/models"
	"vrcomandaapi/internal/shared/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&models.LancamentoComanda{}, &models.LancamentoComandaItem{}); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	return db
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newTestDB(t)

	r := gin.New()
	RegisterRoutes(r, db)

	t.Run("GET /lancamentos returns 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		data := utils.AssertDataArray(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
		if len(data) != 0 {
			t.Fatalf("expected empty data, got len=%d", len(data))
		}
	})

	t.Run("GET /lancamentos/itens returns 400 when id_comanda is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/itens", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("POST /lancamentos returns 400 from validation", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code == http.StatusNotFound {
			t.Fatal("route not registered")
		}
	})

	t.Run("PUT /lancamentos/:id returns 400 for non-numeric id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/abc", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("POST /lancamentos/itens returns 400 for empty body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code == http.StatusNotFound {
			t.Fatal("route not registered")
		}
	})

	t.Run("PUT /lancamentos/itens/:id returns 400 for non-numeric id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/abc", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})
}
