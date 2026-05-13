package parametros

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
	listFn func(ctx context.Context, req ListParametrosRequest) ([]ParametroResponse, error)
}

func (s serviceStub) List(ctx context.Context, req ListParametrosRequest) ([]ParametroResponse, error) {
	return s.listFn(ctx, req)
}

func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, _ ListParametrosRequest) ([]ParametroResponse, error) {
			t.Fatal("service should not be called on bind error")
			return nil, nil
		}})

		r := gin.New()
		r.GET("/parametros", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parametros?idLoja=abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessagePresent(t, body)
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 400 when service rejects request", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListParametrosRequest) ([]ParametroResponse, error) {
			if req.IDLoja != 0 {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return nil, ErrInvalidRequest
		}})

		r := gin.New()
		r.GET("/parametros", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parametros?idLoja=0", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 when service fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListParametrosRequest) ([]ParametroResponse, error) {
			if req.IDLoja != 5 {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return nil, errors.New("boom")
		}})

		r := gin.New()
		r.GET("/parametros", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parametros?idLoja=5", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 with data", func(t *testing.T) {
		called := 0
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListParametrosRequest) ([]ParametroResponse, error) {
			called++
			if req.IDLoja != 7 {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return []ParametroResponse{{IDParametro: 13, Descricao: "Tipo de Etiqueta de Balanca", IDLoja: 7, Valor: "0"}}, nil
		}})

		r := gin.New()
		r.GET("/parametros", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parametros?idLoja=7", nil)
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
