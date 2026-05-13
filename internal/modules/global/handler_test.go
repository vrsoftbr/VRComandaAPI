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

type consultarSituacaoComandaServiceStub struct {
	executeFn func(ctx context.Context, req ConsultarSituacaoComandaRequest) (*ConsultarSituacaoComandaResponse, error)
}

func (s consultarSituacaoComandaServiceStub) Execute(ctx context.Context, req ConsultarSituacaoComandaRequest) (*ConsultarSituacaoComandaResponse, error) {
	if s.executeFn != nil {
		return s.executeFn(ctx, req)
	}
	return nil, nil
}

func newGlobalRouter(service LancamentosDetalhesService, consultarSituacaoService ConsultarSituacaoComandaService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service, consultarSituacaoService)
	r.GET("/lancamentos/detalhes", h.GetLancamentosDetalhes)
	r.GET("/comanda/consultarsituacao", h.ConsultarSituacaoComanda)
	return r
}

func TestHandlerGetLancamentosDetalhes(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarSituacaoComandaServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/detalhes?idLoja=abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when use case returns invalid filter", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{executeFn: func(_ context.Context, _ ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error) {
			return nil, lancamento.ErrInvalidFilter
		}}, consultarSituacaoComandaServiceStub{})

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
		}}, consultarSituacaoComandaServiceStub{})

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
		}}, consultarSituacaoComandaServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/detalhes?idLoja=2&finalizado=true", nil)
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

func TestHandlerConsultarSituacaoComanda(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarSituacaoComandaServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comanda/consultarsituacao?idLoja=abc&numeroIdentificacaoComanda=100", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when use case returns invalid request", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarSituacaoComandaServiceStub{executeFn: func(_ context.Context, _ ConsultarSituacaoComandaRequest) (*ConsultarSituacaoComandaResponse, error) {
			return nil, ErrInvalidRequest
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comanda/consultarsituacao?idLoja=1&numeroIdentificacaoComanda=100", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 with not found message", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarSituacaoComandaServiceStub{executeFn: func(_ context.Context, _ ConsultarSituacaoComandaRequest) (*ConsultarSituacaoComandaResponse, error) {
			return nil, nil
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comanda/consultarsituacao?idLoja=1&numeroIdentificacaoComanda=100", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "Comanda não encontrada")
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 200 with found payload", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarSituacaoComandaServiceStub{executeFn: func(_ context.Context, req ConsultarSituacaoComandaRequest) (*ConsultarSituacaoComandaResponse, error) {
			if req.IDLoja != 1 || req.NumeroIdentificacaoComanda != "1000000001159" {
				t.Fatalf("unexpected req: %+v", req)
			}
			return &ConsultarSituacaoComandaResponse{
				IDLoja:                     1,
				Comanda:                    115,
				NumeroIdentificacaoComanda: "1000000001159",
				Situacao:                   SituacaoComandaBloqueada,
			}, nil
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comanda/consultarsituacao?idLoja=1&numeroIdentificacaoComanda=1000000001159", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "Comanda encontrada")
		data, ok := body["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected object data, got %T", body["data"])
		}
		if data["situacao"] != float64(SituacaoComandaBloqueada) {
			t.Fatalf("unexpected situacao: %v", data["situacao"])
		}
	})
}
