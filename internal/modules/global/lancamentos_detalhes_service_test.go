package global

import (
	"context"
	"errors"
	"testing"
	"time"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/shared/models"
)

type lancamentoServiceStub struct {
	listFn func(ctx context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error)
}

func (s lancamentoServiceStub) Create(ctx context.Context, req lancamento.CreateLancamentoRequest) (*models.LancamentoComanda, error) {
	return nil, nil
}

func (s lancamentoServiceStub) Update(ctx context.Context, id uint, req lancamento.CreateLancamentoRequest) (*models.LancamentoComanda, error) {
	return nil, nil
}

func (s lancamentoServiceStub) List(ctx context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
	if s.listFn != nil {
		return s.listFn(ctx, req)
	}
	return []models.LancamentoComanda{}, nil
}

func (s lancamentoServiceStub) ListItens(ctx context.Context, req lancamento.ListItensRequest) ([]lancamento.ItemComandaResponse, error) {
	return nil, nil
}

func (s lancamentoServiceStub) CreateItems(ctx context.Context, req lancamento.CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
	return nil, nil
}

func (s lancamentoServiceStub) UpdateItem(ctx context.Context, id uint, req lancamento.UpdateItemRequest) (*models.LancamentoComandaItem, error) {
	return nil, nil
}

type comandaServiceStub struct {
	listFn func(ctx context.Context, req comanda.ListComandasRequest) ([]comanda.ComandaResponse, error)
}

func (s comandaServiceStub) List(ctx context.Context, req comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
	if s.listFn != nil {
		return s.listFn(ctx, req)
	}
	return []comanda.ComandaResponse{}, nil
}

type mesaServiceStub struct {
	listFn func(ctx context.Context, req mesa.ListMesasRequest) ([]mesa.MesaResponse, error)
}

func (s mesaServiceStub) List(ctx context.Context, req mesa.ListMesasRequest) ([]mesa.MesaResponse, error) {
	if s.listFn != nil {
		return s.listFn(ctx, req)
	}
	return []mesa.MesaResponse{}, nil
}

func TestLancamentosDetalhesServiceExecute(t *testing.T) {
	t.Run("returns lancamento error", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return nil, errors.New("boom")
			}},
			comandaServiceStub{},
			mesaServiceStub{},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns empty list", func(t *testing.T) {
		service := NewLancamentosDetalhesService(lancamentoServiceStub{}, comandaServiceStub{}, mesaServiceStub{})

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("expected empty result, got %d", len(result))
		}
	})

	t.Run("defaults finalizado to false", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				if req.Finalizado == nil || *req.Finalizado {
					t.Fatalf("expected finalizado=false by default, got %+v", req)
				}
				return []models.LancamentoComanda{}, nil
			}},
			comandaServiceStub{},
			mesaServiceStub{},
		)

		if _, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns composed response", func(t *testing.T) {
		mesaID := 10
		now := time.Date(2026, time.April, 20, 12, 0, 0, 0, time.UTC)

		var capturedLancamentoReq lancamento.ListLancamentosRequest
		var capturedComandaReq comanda.ListComandasRequest
		var capturedMesaReq mesa.ListMesasRequest

		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				capturedLancamentoReq = req
				if req.Finalizado == nil || !*req.Finalizado {
					t.Fatalf("expected finalizado=true, got %+v", req)
				}
				if req.IDLoja != 2 {
					t.Fatalf("expected id_loja=2, got %+v", req)
				}
				return []models.LancamentoComanda{{
					ID:          1,
					IDLoja:      2,
					IDComanda:   100,
					IDMesa:      &mesaID,
					IDAtendente: 7,
					DataHora:    now,
					Finalizado:  false,
					Itens:       []models.LancamentoComandaItem{{ID: 9, IDProduto: 20}},
				}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, req comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				capturedComandaReq = req
				return []comanda.ComandaResponse{{Comanda: 100, NumeroIdentificacao: "A1"}}, nil
			}},
			mesaServiceStub{listFn: func(_ context.Context, req mesa.ListMesasRequest) ([]mesa.MesaResponse, error) {
				capturedMesaReq = req
				return []mesa.MesaResponse{{Mesa: 10, Descricao: "Janela"}}, nil
			}},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{GlobalFilterRequest: GlobalFilterRequest{IDLoja: 2, Finalizado: boolPtr(true)}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedLancamentoReq.Finalizado == nil || !*capturedLancamentoReq.Finalizado || capturedLancamentoReq.IDLoja != 2 {
			t.Fatalf("unexpected lancamentos request: %+v", capturedLancamentoReq)
		}
		if capturedComandaReq.IDLoja != 2 || len(capturedComandaReq.Comandas) != 1 || capturedComandaReq.Comandas[0] != 100 {
			t.Fatalf("unexpected comandas request: %+v", capturedComandaReq)
		}
		if capturedMesaReq.IDLoja != 2 || len(capturedMesaReq.Mesas) != 1 || capturedMesaReq.Mesas[0] != 10 {
			t.Fatalf("unexpected mesas request: %+v", capturedMesaReq)
		}
		if len(result) != 1 || result[0].IDLoja != 2 || result[0].Comanda == nil || result[0].Mesa == nil || len(result[0].Itens) != 1 {
			t.Fatalf("unexpected result: %+v", result)
		}
	})
}

func boolPtr(v bool) *bool { return &v }
