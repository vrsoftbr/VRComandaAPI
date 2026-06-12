package global

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

type comandaCatracaServiceStub struct {
	executeFn func(ctx context.Context, req ComandaCatracaRequest) (*ComandaCatracaResponse, error)
}

func (s comandaCatracaServiceStub) Execute(ctx context.Context, req ComandaCatracaRequest) (*ComandaCatracaResponse, error) {
	if s.executeFn != nil {
		return s.executeFn(ctx, req)
	}
	return nil, nil
}

type comandaPDVServiceStub struct {
	getFn func(ctx context.Context, req GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error)

	updateFn func(ctx context.Context, req UpdadeLancamentoPDVRequest) error
}

func (s comandaPDVServiceStub) Consultar(ctx context.Context, req GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error) {
	if s.getFn != nil {
		return s.getFn(ctx, req)
	}

	return nil, nil
}

func (s comandaPDVServiceStub) Atualizar(ctx context.Context, req UpdadeLancamentoPDVRequest) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, req)
	}

	return nil
}

func newGlobalRouter(service LancamentosDetalhesService, consultarSituacaoService ComandaCatracaService, comandaPDVService ComandaPDVService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(service, consultarSituacaoService, comandaPDVService)
	r.GET("/lancamentos/detalhes", h.GetLancamentosDetalhes)
	r.GET("/comanda/consultarsituacao", h.ComandaCatraca)
	r.GET("/venda/comanda/pdv/consultar", h.GetLancamentosPDV)
	r.PUT("/atualizacomanda", h.UpdateComandaPDV)
	return r
}

func TestHandlerGetLancamentosDetalhes(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, comandaCatracaServiceStub{}, comandaPDVServiceStub{})

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
		}}, comandaCatracaServiceStub{}, comandaPDVServiceStub{})

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
		}}, comandaCatracaServiceStub{}, comandaPDVServiceStub{})

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
		}}, comandaCatracaServiceStub{}, comandaPDVServiceStub{})

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

func TestHandlerComandaCatraca(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, comandaCatracaServiceStub{}, comandaPDVServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comanda/consultarsituacao?idLoja=abc&numeroIdentificacaoComanda=100", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when use case returns invalid request", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, comandaCatracaServiceStub{executeFn: func(_ context.Context, _ ComandaCatracaRequest) (*ComandaCatracaResponse, error) {
			return nil, ErrInvalidRequest
		}}, comandaPDVServiceStub{})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/comanda/consultarsituacao?idLoja=1&numeroIdentificacaoComanda=100", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 with not found message", func(t *testing.T) {
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, comandaCatracaServiceStub{executeFn: func(_ context.Context, _ ComandaCatracaRequest) (*ComandaCatracaResponse, error) {
			return nil, nil
		}}, comandaPDVServiceStub{})

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
		r := newGlobalRouter(lancamentosDetalhesServiceStub{}, comandaCatracaServiceStub{executeFn: func(_ context.Context, req ComandaCatracaRequest) (*ComandaCatracaResponse, error) {
			if req.IDLoja != 1 || req.NumeroIdentificacaoComanda != "1000000001159" {
				t.Fatalf("unexpected req: %+v", req)
			}
			return &ComandaCatracaResponse{
				IDLoja:                     1,
				Comanda:                    115,
				NumeroIdentificacaoComanda: "1000000001159",
				Situacao:                   SituacaoComandaBloqueada,
			}, nil
		}}, comandaPDVServiceStub{})

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

// /
func TestHandlerGetLancamentosPDV(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/venda/comanda/pdv/consultar?numeroComanda=abc&loja=1",
			nil,
		)

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when service returns invalid request", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				getFn: func(_ context.Context, _ GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error) {
					return nil, ErrInvalidRequest
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/venda/comanda/pdv/consultar?numeroComanda=2&loja=1",
			nil,
		)

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when service returns invalid filter", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				getFn: func(_ context.Context, _ GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error) {
					return nil, lancamento.ErrInvalidFilter
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/venda/comanda/pdv/consultar?numeroComanda=2&loja=1",
			nil,
		)

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 when service fails", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				getFn: func(_ context.Context, _ GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error) {
					return nil, errors.New("boom")
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/venda/comanda/pdv/consultar?numeroComanda=2&loja=1",
			nil,
		)

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 with payload", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				getFn: func(_ context.Context, req GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error) {
					if req.NumeroComanda != 2 || req.IDLoja != 1 {
						t.Fatalf("unexpected req: %+v", req)
					}

					return &GetLancamentoItemPDVResponse{
						CodigoComanda:        2,
						TipoDocumentoCliente: 1,
						CodigoVendedor:       123456,
						Itens: []GetLancamentoPDVItemDTO{
							{
								CodigoBarras: "7891000100103",
								Quantidade:   10,
								PrecoVenda:   16.28,
							},
						},
					}, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/venda/comanda/pdv/consultar?numeroComanda=2&loja=1",
			nil,
		)

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

		if data["codigoComanda"] != float64(2) {
			t.Fatalf("unexpected data: %+v", data)
		}
	})
}

func TestHandlerUpdateComandaPDV(t *testing.T) {
	t.Run("returns 400 when body is malformed", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(
			http.MethodPut,
			"/atualizacomanda",
			bytes.NewBufferString("{invalid}"),
		)

		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when required field is missing", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(
			http.MethodPut,
			"/atualizacomanda",
			bytes.NewBufferString(`{"id_loja":1,"id_comanda":[2]}`),
		)

		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 404 when service returns gorm record not found", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				updateFn: func(_ context.Context, _ UpdadeLancamentoPDVRequest) error {
					return gorm.ErrRecordNotFound
				},
			},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(
			http.MethodPut,
			"/atualizacomanda",
			bytes.NewBufferString(`{"id_loja":1,"id_comanda":[2],"finalizado":true}`),
		)

		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 404 when service returns not found", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				updateFn: func(_ context.Context, _ UpdadeLancamentoPDVRequest) error {
					return lancamento.ErrNotFound
				},
			},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(
			http.MethodPut,
			"/atualizacomanda",
			bytes.NewBufferString(`{"id_loja":1,"id_comanda":[2],"finalizado":true}`),
		)

		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when service returns validation error", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				updateFn: func(_ context.Context, _ UpdadeLancamentoPDVRequest) error {
					return lancamento.ErrValidation
				},
			},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(
			http.MethodPut,
			"/atualizacomanda",
			bytes.NewBufferString(`{"id_loja":1,"id_comanda":[2],"finalizado":true}`),
		)

		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 when service fails", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				updateFn: func(_ context.Context, _ UpdadeLancamentoPDVRequest) error {
					return errors.New("boom")
				},
			},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(
			http.MethodPut,
			"/atualizacomanda",
			bytes.NewBufferString(`{"id_loja":1,"id_comanda":[2],"finalizado":true}`),
		)

		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 on success", func(t *testing.T) {
		r := newGlobalRouter(
			lancamentosDetalhesServiceStub{},
			comandaCatracaServiceStub{},
			comandaPDVServiceStub{
				updateFn: func(_ context.Context, req UpdadeLancamentoPDVRequest) error {
					if req.IDLoja != 1 {
						t.Fatalf("unexpected loja")
					}

					if len(req.IDComanda) != 1 || req.IDComanda[0] != 2 {
						t.Fatalf("unexpected comandas: %+v", req.IDComanda)
					}

					if req.Finalizado == nil || !*req.Finalizado {
						t.Fatalf("unexpected finalizado")
					}

					return nil
				},
			},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(
			http.MethodPut,
			"/atualizacomanda",
			bytes.NewBufferString(`{"id_loja":1,"id_comanda":[2],"finalizado":true}`),
		)

		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})
}
