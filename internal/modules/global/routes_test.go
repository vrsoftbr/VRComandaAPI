package global

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, lancamentoServiceStub{}, comandaServiceStub{}, mesaServiceStub{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/lancamentos/detalhes", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected registered route status=200, got %d", w.Code)
	}
}
