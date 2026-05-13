package produto

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"vrcomandaapi/internal/shared/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("wires produtos endpoint", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		r := gin.New()
		RegisterRoutes(r, func() *mongo.Database { return nil }, func() {})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/produtos", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "ok")
		data, ok := body["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object, got %T", body["data"])
		}
		items, ok := data["items"].([]any)
		if !ok {
			t.Fatalf("expected items array, got %T", data["items"])
		}
		if len(items) != 0 {
			t.Fatalf("expected empty items array, got len=%d", len(items))
		}
	})
}
