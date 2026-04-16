package comanda

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
	listFn func(ctx context.Context, req ListComandasRequest) ([]ComandaResponse, error)
}

func (s serviceStub) List(ctx context.Context, req ListComandasRequest) ([]ComandaResponse, error) {
	return s.listFn(ctx, req)
}

func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, _ ListComandasRequest) ([]ComandaResponse, error) {
			t.Fatal("service should not be called on bind error")
			return nil, nil
		}})

		r := gin.New()
		r.GET("/comandas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comandas?ativo=not-bool", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessagePresent(t, body)
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 500 when service fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListComandasRequest) ([]ComandaResponse, error) {
			if req.IDLoja != 5 || req.Comanda != 10 || req.NumeroIdentificacao != "AA" {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return nil, errors.New("boom")
		}})

		r := gin.New()
		r.GET("/comandas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comandas?idLoja=5&comanda=10&numeroIdentificacao=AA", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 200 with data", func(t *testing.T) {
		called := 0
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListComandasRequest) ([]ComandaResponse, error) {
			called++
			if req.IDLoja != 7 || req.Comanda != 70 || req.NumeroIdentificacao != "ID70" || req.Ativo == nil || !*req.Ativo {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return []ComandaResponse{{ID: "1", IDLoja: 7, Comanda: 70, NumeroIdentificacao: "ID70", Observacao: "obs", Ativo: true}}, nil
		}})

		r := gin.New()
		r.GET("/comandas", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comandas?idLoja=7&comanda=70&numeroIdentificacao=ID70&ativo=true", nil)
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
