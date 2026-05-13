package global

import (
	"context"
	"errors"
	"testing"
	"time"

	"vrcomandaapi/internal/modules/atendente"
	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/modules/produto"
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

type atendenteServiceStub struct {
	listFn func(ctx context.Context, req atendente.ListAtendentesRequest) ([]atendente.AtendenteResponse, error)
}

func (s atendenteServiceStub) List(ctx context.Context, req atendente.ListAtendentesRequest) ([]atendente.AtendenteResponse, error) {
	if s.listFn != nil {
		return s.listFn(ctx, req)
	}
	return []atendente.AtendenteResponse{}, nil
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

type produtoServiceStub struct {
	listFn func(ctx context.Context, req produto.ListProdutosRequest) (interface{}, error)
}

func (s produtoServiceStub) List(ctx context.Context, req produto.ListProdutosRequest) (interface{}, error) {
	if s.listFn != nil {
		return s.listFn(ctx, req)
	}
	return produto.ProdutosPaginatedResponse{Items: []produto.ProdutoResponse{}, Page: 1, Limit: 20, Total: 0, Pages: 0}, nil
}

func TestLancamentosDetalhesServiceExecute(t *testing.T) {
	t.Run("returns lancamento error", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return nil, errors.New("boom")
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns empty list", func(t *testing.T) {
		service := NewLancamentosDetalhesService(lancamentoServiceStub{}, atendenteServiceStub{}, comandaServiceStub{}, mesaServiceStub{}, produtoServiceStub{})

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
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		if _, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns composed response", func(t *testing.T) {
		mesaID := 10
		now := time.Date(2026, time.April, 20, 12, 0, 0, 0, time.UTC)

		var capturedLancamentoReq lancamento.ListLancamentosRequest
		var capturedAtendenteReq atendente.ListAtendentesRequest
		var capturedComandaReq comanda.ListComandasRequest
		var capturedMesaReq mesa.ListMesasRequest
		var capturedProdutoReq produto.ListProdutosRequest

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
			atendenteServiceStub{listFn: func(_ context.Context, req atendente.ListAtendentesRequest) ([]atendente.AtendenteResponse, error) {
				capturedAtendenteReq = req
				return []atendente.AtendenteResponse{{IDLoja: req.IDLoja, IDAtendente: req.IDAtendente, Nome: "Maria"}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, req comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				capturedComandaReq = req
				return []comanda.ComandaResponse{{Comanda: 100, NumeroIdentificacao: "A1"}}, nil
			}},
			mesaServiceStub{listFn: func(_ context.Context, req mesa.ListMesasRequest) ([]mesa.MesaResponse, error) {
				capturedMesaReq = req
				return []mesa.MesaResponse{{Mesa: 10, Descricao: "Janela"}}, nil
			}},
			produtoServiceStub{listFn: func(_ context.Context, req produto.ListProdutosRequest) (interface{}, error) {
				capturedProdutoReq = req
				return produto.ProdutosPaginatedResponse{Items: []produto.ProdutoResponse{{CodigoBarras: req.CodigoBarras, DescricaoCompleta: "Refrigerante"}}, Page: 1, Limit: 20, Total: 1, Pages: 1}, nil
			}},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{GlobalFilterRequest: GlobalFilterRequest{IDLoja: 2, Finalizado: boolPtr(true)}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedLancamentoReq.Finalizado == nil || !*capturedLancamentoReq.Finalizado || capturedLancamentoReq.IDLoja != 2 {
			t.Fatalf("unexpected lancamentos request: %+v", capturedLancamentoReq)
		}
		if capturedAtendenteReq.IDLoja != 2 || capturedAtendenteReq.IDAtendente != "7" {
			t.Fatalf("unexpected atendentes request: %+v", capturedAtendenteReq)
		}
		if capturedComandaReq.IDLoja != 2 || len(capturedComandaReq.Comandas) != 1 || capturedComandaReq.Comandas[0] != 100 {
			t.Fatalf("unexpected comandas request: %+v", capturedComandaReq)
		}
		if capturedMesaReq.IDLoja != 2 || len(capturedMesaReq.Mesas) != 1 || capturedMesaReq.Mesas[0] != 10 {
			t.Fatalf("unexpected mesas request: %+v", capturedMesaReq)
		}
		if capturedProdutoReq.IDLoja != 0 || capturedProdutoReq.CodigoBarras != "" {
			t.Fatalf("did not expect produtos query when barcode is empty, got %+v", capturedProdutoReq)
		}
		if len(result) != 1 || result[0].IDLoja != 2 || result[0].Comanda == nil || result[0].Mesa == nil || len(result[0].Itens) != 1 {
			t.Fatalf("unexpected result: %+v", result)
		}
		if result[0].Atendente == nil || result[0].Atendente.Nome != "Maria" {
			t.Fatalf("expected atendente to be filled, got %+v", result[0].Atendente)
		}
		if result[0].Itens[0].DescricaoProduto != "" {
			t.Fatalf("expected empty product name for item without barcode, got %+v", result[0].Itens[0])
		}
	})

	t.Run("matches atendente by lancamento store and attendant id", func(t *testing.T) {
		capturedAtendenteReqs := make([]atendente.ListAtendentesRequest, 0)

		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{
					{ID: 1, IDLoja: 2, IDComanda: 100, IDAtendente: 7, DataHora: time.Date(2026, time.May, 10, 10, 0, 0, 0, time.UTC)},
					{ID: 2, IDLoja: 3, IDComanda: 101, IDAtendente: 7, DataHora: time.Date(2026, time.May, 10, 11, 0, 0, 0, time.UTC)},
				}, nil
			}},
			atendenteServiceStub{listFn: func(_ context.Context, req atendente.ListAtendentesRequest) ([]atendente.AtendenteResponse, error) {
				capturedAtendenteReqs = append(capturedAtendenteReqs, req)
				return []atendente.AtendenteResponse{{IDLoja: req.IDLoja, IDAtendente: req.IDAtendente, Nome: "Atendente Loja " + req.IDAtendente + "-" + string(rune('0'+req.IDLoja))}}, nil
			}},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{GlobalFilterRequest: GlobalFilterRequest{IDLoja: 2}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(capturedAtendenteReqs) != 2 {
			t.Fatalf("expected two atendente queries, got %+v", capturedAtendenteReqs)
		}
		if capturedAtendenteReqs[0].IDLoja != 2 || capturedAtendenteReqs[0].IDAtendente != "7" {
			t.Fatalf("unexpected first atendente request: %+v", capturedAtendenteReqs[0])
		}
		if capturedAtendenteReqs[1].IDLoja != 3 || capturedAtendenteReqs[1].IDAtendente != "7" {
			t.Fatalf("unexpected second atendente request: %+v", capturedAtendenteReqs[1])
		}
		if len(result) != 2 {
			t.Fatalf("unexpected result length: %+v", result)
		}
		if result[0].Atendente == nil || result[0].Atendente.IDLoja != 2 || result[0].Atendente.IDAtendente != "7" {
			t.Fatalf("unexpected first atendente result: %+v", result[0].Atendente)
		}
		if result[1].Atendente == nil || result[1].Atendente.IDLoja != 3 || result[1].Atendente.IDAtendente != "7" {
			t.Fatalf("unexpected second atendente result: %+v", result[1].Atendente)
		}
	})

	t.Run("includes product name in item by barcode and store", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{
					ID:        1,
					IDLoja:    5,
					IDComanda: 42,
					DataHora:  time.Date(2026, time.January, 2, 10, 0, 0, 0, time.UTC),
					Itens: []models.LancamentoComandaItem{{
						ID:           10,
						CodigoBarras: "789",
					}},
				}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{listFn: func(_ context.Context, req produto.ListProdutosRequest) (interface{}, error) {
				if req.IDLoja != 5 || req.CodigoBarras != "789" {
					t.Fatalf("unexpected produto request: %+v", req)
				}
				return produto.ProdutosPaginatedResponse{Items: []produto.ProdutoResponse{{CodigoBarras: "789", DescricaoCompleta: "Coca-Cola 350ml"}}, Page: 1, Limit: 20, Total: 1, Pages: 1}, nil
			}},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{GlobalFilterRequest: GlobalFilterRequest{IDLoja: 5}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || len(result[0].Itens) != 1 {
			t.Fatalf("unexpected result: %+v", result)
		}
		if result[0].Itens[0].DescricaoProduto != "Coca-Cola 350ml" {
			t.Fatalf("expected nome_produto to be filled, got %+v", result[0].Itens[0])
		}
	})

	t.Run("propagates atendente service error", func(t *testing.T) {
		expectedErr := errors.New("atendente boom")
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDLoja: 1, IDComanda: 10, IDAtendente: 5}}, nil
			}},
			atendenteServiceStub{listFn: func(_ context.Context, _ atendente.ListAtendentesRequest) ([]atendente.AtendenteResponse, error) {
				return nil, expectedErr
			}},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("propagates comanda service error", func(t *testing.T) {
		expectedErr := errors.New("comanda boom")
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDLoja: 1, IDComanda: 10}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return nil, expectedErr
			}},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("propagates mesa service error", func(t *testing.T) {
		expectedErr := errors.New("mesa boom")
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDLoja: 1, IDComanda: 10}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{listFn: func(_ context.Context, _ mesa.ListMesasRequest) ([]mesa.MesaResponse, error) {
				return nil, expectedErr
			}},
			produtoServiceStub{},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("propagates produto service error from listProdutosByBarcode", func(t *testing.T) {
		expectedErr := errors.New("produto boom")
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{
					ID:        1,
					IDLoja:    1,
					IDComanda: 10,
					Itens:     []models.LancamentoComandaItem{{CodigoBarras: "123"}},
				}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{listFn: func(_ context.Context, _ produto.ListProdutosRequest) (interface{}, error) {
				return nil, expectedErr
			}},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("handles atendente not found gracefully", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDLoja: 1, IDComanda: 10, IDAtendente: 99}}, nil
			}},
			atendenteServiceStub{listFn: func(_ context.Context, _ atendente.ListAtendentesRequest) ([]atendente.AtendenteResponse, error) {
				return []atendente.AtendenteResponse{}, nil // not found
			}},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || result[0].Atendente != nil {
			t.Fatalf("expected one result with nil Atendente, got %+v", result)
		}
	})

	t.Run("handles duplicate lancamento comanda IDs", func(t *testing.T) {
		var capturedComandasReq comanda.ListComandasRequest
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{
					{ID: 1, IDLoja: 1, IDComanda: 10},
					{ID: 2, IDLoja: 1, IDComanda: 10}, // duplicate comanda ID
					{ID: 3, IDLoja: 1, IDComanda: 20},
				}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{listFn: func(_ context.Context, req comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				capturedComandasReq = req
				return []comanda.ComandaResponse{}, nil
			}},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(capturedComandasReq.Comandas) != 2 || capturedComandasReq.Comandas[0] != 10 || capturedComandasReq.Comandas[1] != 20 {
			t.Fatalf("expected deduplicated sorted comandas [10, 20], got %+v", capturedComandasReq.Comandas)
		}
	})

	t.Run("handles duplicate atendente keys across lancamentos", func(t *testing.T) {
		atendenteCalls := 0
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{
					{ID: 1, IDLoja: 1, IDComanda: 10, IDAtendente: 5},
					{ID: 2, IDLoja: 1, IDComanda: 11, IDAtendente: 5}, // same atendente
				}, nil
			}},
			atendenteServiceStub{listFn: func(_ context.Context, req atendente.ListAtendentesRequest) ([]atendente.AtendenteResponse, error) {
				atendenteCalls++
				return []atendente.AtendenteResponse{{IDLoja: req.IDLoja, IDAtendente: "5", Nome: "Ana"}}, nil
			}},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if atendenteCalls != 1 {
			t.Fatalf("expected 1 atendente query for deduped key, got %d", atendenteCalls)
		}
	})

	t.Run("listProdutosByBarcode handles pointer paginated response type", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{
					ID:        1,
					IDLoja:    1,
					IDComanda: 10,
					Itens:     []models.LancamentoComandaItem{{CodigoBarras: "456"}},
				}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{listFn: func(_ context.Context, _ produto.ListProdutosRequest) (interface{}, error) {
				r := &produto.ProdutosPaginatedResponse{
					Items: []produto.ProdutoResponse{{CodigoBarras: "456", DescricaoCompleta: "Item Pointer"}},
				}
				return r, nil
			}},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || result[0].Itens[0].DescricaoProduto != "Item Pointer" {
			t.Fatalf("expected DescricaoProduto 'Item Pointer', got %+v", result)
		}
	})

	t.Run("listProdutosByBarcode handles nil pointer paginated response", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{
					ID:        1,
					IDLoja:    1,
					IDComanda: 10,
					Itens:     []models.LancamentoComandaItem{{CodigoBarras: "789"}},
				}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{listFn: func(_ context.Context, _ produto.ListProdutosRequest) (interface{}, error) {
				var r *produto.ProdutosPaginatedResponse
				return r, nil // nil pointer
			}},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || result[0].Itens[0].DescricaoProduto != "" {
			t.Fatalf("expected empty DescricaoProduto for nil pointer response, got %+v", result)
		}
	})

	t.Run("listProdutosByBarcode handles slice response type", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{
					ID:        1,
					IDLoja:    1,
					IDComanda: 10,
					Itens:     []models.LancamentoComandaItem{{CodigoBarras: "AAA"}},
				}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{listFn: func(_ context.Context, _ produto.ListProdutosRequest) (interface{}, error) {
				return []produto.ProdutoResponse{{CodigoBarras: "AAA", DescricaoCompleta: "Slice Item"}}, nil
			}},
		)

		result, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || result[0].Itens[0].DescricaoProduto != "Slice Item" {
			t.Fatalf("expected DescricaoProduto 'Slice Item', got %+v", result)
		}
	})

	t.Run("listProdutosByBarcode returns error for unexpected response type", func(t *testing.T) {
		service := NewLancamentosDetalhesService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{
					ID:        1,
					IDLoja:    1,
					IDComanda: 10,
					Itens:     []models.LancamentoComandaItem{{CodigoBarras: "BBB"}},
				}}, nil
			}},
			atendenteServiceStub{},
			comandaServiceStub{},
			mesaServiceStub{},
			produtoServiceStub{listFn: func(_ context.Context, _ produto.ListProdutosRequest) (interface{}, error) {
				return "unexpected string type", nil
			}},
		)

		_, err := service.Execute(context.Background(), ListLancamentosDetalhesRequest{})
		if err == nil {
			t.Fatal("expected error for unexpected produto response type, got nil")
		}
	})

	t.Run("buildAtendentesByKey skips non-numeric IDAtendente", func(t *testing.T) {
		items := []atendente.AtendenteResponse{
			{IDLoja: 1, IDAtendente: "abc", Nome: "Ana"},  // non-numeric, should be skipped
			{IDLoja: 1, IDAtendente: "5", Nome: "Carlos"}, // valid
		}
		result := buildAtendentesByKey(items)
		if len(result) != 1 {
			t.Fatalf("expected 1 entry (skipped non-numeric), got %d entries: %+v", len(result), result)
		}
		if _, ok := result["1:5"]; !ok {
			t.Fatalf("expected key '1:5' in result, got %+v", result)
		}
	})

	t.Run("buildProdutosByBarcode skips empty CodigoBarras", func(t *testing.T) {
		items := []produto.ProdutoResponse{
			{CodigoBarras: "", DescricaoCompleta: "No barcode"},
			{CodigoBarras: "  ", DescricaoCompleta: "Whitespace barcode"},
			{CodigoBarras: "123", DescricaoCompleta: "Valid"},
		}
		result := buildProdutosByBarcode(items)
		if len(result) != 1 {
			t.Fatalf("expected 1 entry, got %d: %+v", len(result), result)
		}
		if _, ok := result["123"]; !ok {
			t.Fatalf("expected key '123', got %+v", result)
		}
	})

	t.Run("uniqueSortedItemBarcodes skips empty CodigoBarras", func(t *testing.T) {
		lancamentos := []models.LancamentoComanda{{
			ID: 1,
			Itens: []models.LancamentoComandaItem{
				{CodigoBarras: ""},
				{CodigoBarras: "  "},
				{CodigoBarras: "abc"},
				{CodigoBarras: "abc"}, // duplicate
			},
		}}
		result := uniqueSortedItemBarcodes(lancamentos)
		if len(result) != 1 || result[0] != "abc" {
			t.Fatalf("expected [abc], got %+v", result)
		}
	})

	t.Run("uniqueSortedMesaIDs skips nil IDMesa", func(t *testing.T) {
		mesaID := 5
		lancamentos := []models.LancamentoComanda{
			{ID: 1, IDMesa: nil},
			{ID: 2, IDMesa: &mesaID},
			{ID: 3, IDMesa: &mesaID}, // duplicate
		}
		result := uniqueSortedMesaIDs(lancamentos)
		if len(result) != 1 || result[0] != 5 {
			t.Fatalf("expected [5], got %+v", result)
		}
	})

	t.Run("cloneItens returns empty slice when itens is empty", func(t *testing.T) {
		result := cloneItens([]models.LancamentoComandaItem{}, map[string]*produto.ProdutoResponse{})
		if result == nil || len(result) != 0 {
			t.Fatalf("expected empty slice, got %+v", result)
		}
	})

	t.Run("cloneItens falls back to DescricaoCupom when DescricaoCompleta is empty", func(t *testing.T) {
		p := &produto.ProdutoResponse{CodigoBarras: "X1", DescricaoCompleta: "", DescricaoCupom: "Cupom Name"}
		itens := []models.LancamentoComandaItem{{CodigoBarras: "X1"}}
		result := cloneItens(itens, map[string]*produto.ProdutoResponse{"X1": p})
		if len(result) != 1 || result[0].DescricaoProduto != "Cupom Name" {
			t.Fatalf("expected DescricaoProduto 'Cupom Name', got %+v", result)
		}
	})
}

func boolPtr(v bool) *bool { return &v }
