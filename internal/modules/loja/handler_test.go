package loja

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/shared/utils"
)

type serviceStub struct {
	listFn func(ctx context.Context) ([]LojaResponse, error)
}

func (s serviceStub) List(ctx context.Context) ([]LojaResponse, error) {
	return s.listFn(ctx)
}

func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 500 when service fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context) ([]LojaResponse, error) {
			return nil, errors.New("boom")
		}})

		r := gin.New()
		r.GET("/lojas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lojas", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 200 with data", func(t *testing.T) {
		called := 0
		h := NewHandler(serviceStub{listFn: func(_ context.Context) ([]LojaResponse, error) {
			called++
			return []LojaResponse{{ID: 1, Descricao: "Loja 1", NomeFantasia: "Fantasia"}}, nil
		}})

		r := gin.New()
		r.GET("/lojas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lojas", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		if called != 1 {
			t.Fatalf("service calls = %d", called)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "ok")
		data := utils.AssertDataArray(t, body)
		if len(data) != 1 {
			t.Fatalf("expected one item in data, got %d", len(data))
		}
	})
}
