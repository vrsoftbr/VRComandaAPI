package mesa

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"vrcomandaapi/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

type serviceStub struct {
	listFn func(ctx context.Context, req ListMesasRequest) ([]MesaResponse, error)
}

func (s serviceStub) List(ctx context.Context, req ListMesasRequest) ([]MesaResponse, error) {
	return s.listFn(ctx, req)
}

func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, _ ListMesasRequest) ([]MesaResponse, error) {
			t.Fatal("service should not be called on bind error")
			return nil, nil
		}})

		r := gin.New()
		r.GET("/mesas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mesas?ativo=not-bool", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessagePresent(t, body)
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 500 when service fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListMesasRequest) ([]MesaResponse, error) {
			if req.IDLoja != 5 || req.Mesa != 10 {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return nil, errors.New("boom")
		}})

		r := gin.New()
		r.GET("/mesas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mesas?idLoja=5&mesa=10", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 200 with data", func(t *testing.T) {
		called := 0
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListMesasRequest) ([]MesaResponse, error) {
			called++
			if req.IDLoja != 7 || req.Mesa != 70 || req.Ativo == nil || !*req.Ativo {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return []MesaResponse{{ID: "1", IDLoja: 7, Mesa: 70, Descricao: "M70", Ativo: true}}, nil
		}})

		r := gin.New()
		r.GET("/mesas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/mesas?idLoja=7&mesa=70&ativo=true", nil)
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
