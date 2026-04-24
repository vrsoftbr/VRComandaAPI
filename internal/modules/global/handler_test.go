package global

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/utils"
)

type lancamentosDetalhesServiceStub struct {
	executeFn func(ctx context.Context, req ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error)
}

func (s lancamentosDetalhesServiceStub) Execute(ctx context.Context, req ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error) {
	if s.executeFn != nil {
		return s.executeFn(ctx, req)
	}
	return []LancamentoDetalhesDTO{}, nil
}

func newGlobalRouter(service LancamentosDetalhesService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service)
	r.GET("/lancamentos/detalhes", h.GetLancamentosDetalhes)
	return r
}

func TestHandlerGetLancamentosDetalhes(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/detalhes?id_loja=abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when use case returns invalid filter", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{executeFn: func(_ context.Context, _ ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error) {
			return nil, lancamento.ErrInvalidFilter
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/detalhes", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{executeFn: func(_ context.Context, _ ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error) {
			return nil, errors.New("boom")
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/detalhes", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 on success", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{executeFn: func(_ context.Context, req ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error) {
			if req.IDLoja != 2 {
				t.Fatalf("unexpected req: %+v", req)
			}
			return []LancamentoDetalhesDTO{{IDLancamento: 1, IDLoja: 2}}, nil
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/detalhes?id_loja=2&finalizado=true", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		data := utils.AssertDataArray(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
		if len(data) != 1 {
			t.Fatalf("expected one row, got %d", len(data))
		}
	})
}
