package comanda

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"vrcomandaapi/internal/shared/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("wires comandas endpoint", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		r := gin.New()
		RegisterRoutes(r, func() *mongo.Database { return nil }, func() {})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comandas", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "ok")
		data := utils.AssertDataArray(t, body)
		if len(data) != 0 {
			t.Fatalf("expected empty data array, got len=%d", len(data))
		}
	})
}
