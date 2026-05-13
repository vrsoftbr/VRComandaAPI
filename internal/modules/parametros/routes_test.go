package parametros

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"vrcomandaapi/internal/shared/utils"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("wires parametros endpoint", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		r := gin.New()
		RegisterRoutes(r, func() *mongo.Database { return nil }, func() {})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parametros?idLoja=3", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "ok")
		data := utils.AssertDataArray(t, body)
		if len(data) != len(defaultParametros) {
			t.Fatalf("expected %d items, got %d", len(defaultParametros), len(data))
		}
	})
}
