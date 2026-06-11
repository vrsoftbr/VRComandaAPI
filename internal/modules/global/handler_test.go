package global

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/models"
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

type consultarComandaCatracaServiceStub struct {
	executeFn func(ctx context.Context, req ConsultarComandaCatracaRequest) (*ConsultarComandaCatracaResponse, error)
}

func (s consultarComandaCatracaServiceStub) Execute(ctx context.Context, req ConsultarComandaCatracaRequest) (*ConsultarComandaCatracaResponse, error) {
	if s.executeFn != nil {
		return s.executeFn(ctx, req)
	}
	return nil, nil
}

type comandaPDVServiceStub struct {
	consultarFn func(ctx context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error)
	atualizarFn func(ctx context.Context, req AtualizarComandaPDVRequest) (*models.LancamentoComanda, error)
}

func (s comandaPDVServiceStub) Consultar(ctx context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error) {
	if s.consultarFn != nil {
		return s.consultarFn(ctx, req)
	}
	return nil, nil
}

func (s comandaPDVServiceStub) Atualizar(ctx context.Context, req AtualizarComandaPDVRequest) (*models.LancamentoComanda, error) {
	if s.atualizarFn != nil {
		return s.atualizarFn(ctx, req)
	}
	return nil, nil
}

func newGlobalRouter(service LancamentosDetalhesService, consultarSituacaoService ConsultarComandaCatracaService, comandaPDVService ...ComandaPDVService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	pdvService := ComandaPDVService(comandaPDVServiceStub{})
	if len(comandaPDVService) > 0 {
		pdvService = comandaPDVService[0]
	}
	h := NewHandler(service, consultarSituacaoService, pdvService)
	r.GET("/lancamentos/detalhes", h.GetLancamentosDetalhes)
	r.GET("/comanda/consultarsituacao", h.ConsultarComandaCatraca)
	r.GET("/venda/comanda/pdv/consultar", h.ConsultarComandaPDV)
	r.PUT("/atualizacomanda", h.AtualizarComandaPDV)
	return r
}

func TestHandlerGetLancamentosDetalhes(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{})

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
		}}, consultarComandaCatracaServiceStub{})

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
		}}, consultarComandaCatracaServiceStub{})

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
		}}, consultarComandaCatracaServiceStub{})

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

func TestHandlerConsultarComandaCatraca(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comanda/consultarsituacao?idLoja=abc&numeroIdentificacaoComanda=100", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when use case returns invalid request", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{executeFn: func(_ context.Context, _ ConsultarComandaCatracaRequest) (*ConsultarComandaCatracaResponse, error) {
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
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{executeFn: func(_ context.Context, _ ConsultarComandaCatracaRequest) (*ConsultarComandaCatracaResponse, error) {
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
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{executeFn: func(_ context.Context, req ConsultarComandaCatracaRequest) (*ConsultarComandaCatracaResponse, error) {
			if req.IDLoja != 1 || req.NumeroIdentificacaoComanda != "1000000001159" {
				t.Fatalf("unexpected req: %+v", req)
			}
			return &ConsultarComandaCatracaResponse{
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

func TestHandlerConsultarComandaPDV(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/venda/comanda/pdv/consultar?numeroComanda=abc&loja=1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when use case returns invalid request", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{consultarFn: func(_ context.Context, _ ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error) {
				return nil, ErrInvalidRequest
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/venda/comanda/pdv/consultar?numeroComanda=2", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{consultarFn: func(_ context.Context, _ ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error) {
				return nil, errors.New("boom")
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/venda/comanda/pdv/consultar?numeroComanda=2&loja=1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 with fixed message and payload", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{consultarFn: func(_ context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error) {
				if req.NumeroComanda != 2 || req.IDLoja != 1 {
					t.Fatalf("unexpected req: %+v", req)
				}
				return &ConsultarComandaPDVResponse{
					CodigoComanda:        2,
					TipoDocumentoCliente: 1,
					CodigoVendedor:       123456,
					Itens: []ConsultarComandaPDVItemDTO{{
						CodigoBarras: "7891000100103",
						Quantidade:   10,
						PrecoVenda:   16.28,
					}},
				}, nil
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/venda/comanda/pdv/consultar?numeroComanda=2&loja=1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "Comanda")
		data, ok := body["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected object data, got %T", body["data"])
		}
		if data["codigoComanda"] != float64(2) || data["tipoDocumentoCliente"] != float64(1) || data["codigoVendedor"] != float64(123456) {
			t.Fatalf("unexpected data: %+v", data)
		}
	})

	t.Run("accepts idLoja alias query", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{consultarFn: func(_ context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error) {
				if req.IDLoja != 3 {
					t.Fatalf("expected idLoja alias to fill IDLoja, got %+v", req)
				}
				return nil, nil
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/venda/comanda/pdv/consultar?numeroComanda=2&idLoja=3", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})
}

func TestHandlerAtualizarComandaPDV(t *testing.T) {
	t.Run("returns 400 when body is malformed JSON", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/atualizacomanda", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when finalizado is missing", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, consultarComandaCatracaServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/atualizacomanda", bytes.NewBufferString(`{"id_loja":1,"id_comanda":2}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when use case returns invalid request", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{atualizarFn: func(_ context.Context, _ AtualizarComandaPDVRequest) (*models.LancamentoComanda, error) {
				return nil, ErrInvalidRequest
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/atualizacomanda", bytes.NewBufferString(`{"id_loja":1,"id_comanda":2,"finalizado":true}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 404 when use case returns not found", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{atualizarFn: func(_ context.Context, _ AtualizarComandaPDVRequest) (*models.LancamentoComanda, error) {
				return nil, lancamento.ErrNotFound
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/atualizacomanda", bytes.NewBufferString(`{"id_loja":1,"id_comanda":2,"finalizado":true}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{atualizarFn: func(_ context.Context, _ AtualizarComandaPDVRequest) (*models.LancamentoComanda, error) {
				return nil, errors.New("boom")
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/atualizacomanda", bytes.NewBufferString(`{"id_loja":1,"id_comanda":2,"finalizado":true}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 with standard response on success", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			consultarComandaCatracaServiceStub{},
			comandaPDVServiceStub{atualizarFn: func(_ context.Context, req AtualizarComandaPDVRequest) (*models.LancamentoComanda, error) {
				if req.IDLoja != 1 || req.IDComanda != 2 || req.Finalizado == nil || !*req.Finalizado {
					t.Fatalf("unexpected req: %+v", req)
				}
				return &models.LancamentoComanda{ID: 10, IDLoja: req.IDLoja, IDComanda: req.IDComanda, Finalizado: *req.Finalizado}, nil
			}},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/atualizacomanda", bytes.NewBufferString(`{"id_loja":1,"id_comanda":2,"finalizado":true}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "ok")
		data, ok := body["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected object data, got %T", body["data"])
		}
		if data["finalizado"] != true {
			t.Fatalf("unexpected data: %+v", data)
		}
	})
}
