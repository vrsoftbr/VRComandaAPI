package global

import (
	"context"
	"errors"
	"testing"

	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/models"
)

func TestComandaPDVServiceConsultar(t *testing.T) {
	t.Run("returns invalid request when numeroComanda is invalid", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{})

		_, err := svc.Consultar(context.Background(), ConsultarComandaPDVRequest{NumeroComanda: 0, IDLoja: 1})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when loja is invalid", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{})

		_, err := svc.Consultar(context.Background(), ConsultarComandaPDVRequest{NumeroComanda: 2, IDLoja: 0})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("accepts idLoja alias", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{listFn: func(_ context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
			if req.IDLoja != 7 {
				t.Fatalf("expected id_loja=7 from alias, got %+v", req)
			}
			return []models.LancamentoComanda{}, nil
		}})

		_, err := svc.Consultar(context.Background(), ConsultarComandaPDVRequest{NumeroComanda: 2, IDLojaAlias: 7})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns nil when comanda has no open lancamentos", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{})

		res, err := svc.Consultar(context.Background(), ConsultarComandaPDVRequest{NumeroComanda: 2, IDLoja: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})

	t.Run("returns composed pdv payload", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{listFn: func(_ context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
			if req.IDLoja != 1 || req.IDComanda != 2 || req.Finalizado == nil || *req.Finalizado {
				t.Fatalf("unexpected lancamento request: %+v", req)
			}
			return []models.LancamentoComanda{{
				IDLoja:      1,
				IDComanda:   2,
				IDAtendente: 123456,
				Itens: []models.LancamentoComandaItem{
					{CodigoBarras: "7891000100103", Quantidade: 10, PrecoVenda: 16.28},
					{CodigoBarras: "10", Quantidade: 1.568, PrecoVenda: 19.75},
					{CodigoBarras: "999", Quantidade: 1, PrecoVenda: 2, Cancelado: true},
				},
			}}, nil
		}})

		res, err := svc.Consultar(context.Background(), ConsultarComandaPDVRequest{NumeroComanda: 2, IDLoja: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil {
			t.Fatal("expected response, got nil")
		}
		if res.CodigoComanda != 2 || res.TipoDocumentoCliente != 1 || res.DocumentoCliente != "" || res.NomeCliente != "" || res.CodigoVendedor != 123456 {
			t.Fatalf("unexpected header response: %+v", res)
		}
		if res.ValorDescontoVenda != 0 || res.ValorAcrescimoVenda != 0 {
			t.Fatalf("expected zero sale adjustments, got %+v", res)
		}
		if len(res.Itens) != 2 {
			t.Fatalf("expected two non-canceled items, got %+v", res.Itens)
		}
		if res.Itens[0].CodigoBarras != "7891000100103" || res.Itens[0].Quantidade != 10 || res.Itens[0].PrecoVenda != 16.28 {
			t.Fatalf("unexpected first item: %+v", res.Itens[0])
		}
		if res.Itens[1].CodigoBarras != "10" || res.Itens[1].Quantidade != 1.568 || res.Itens[1].PrecoVenda != 19.75 {
			t.Fatalf("unexpected second item: %+v", res.Itens[1])
		}
		if res.Itens[0].ValorDesconto != 0 || res.Itens[0].ValorAcrescimo != 0 {
			t.Fatalf("expected zero item adjustments, got %+v", res.Itens[0])
		}
	})

	t.Run("propagates lancamento service error", func(t *testing.T) {
		expectedErr := errors.New("lancamento failed")
		svc := NewComandaPDVService(lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
			return nil, expectedErr
		}})

		_, err := svc.Consultar(context.Background(), ConsultarComandaPDVRequest{NumeroComanda: 2, IDLoja: 1})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestComandaPDVServiceAtualizar(t *testing.T) {
	t.Run("returns invalid request when id_loja is invalid", func(t *testing.T) {
		finalizado := true
		svc := NewComandaPDVService(lancamentoServiceStub{})

		_, err := svc.Atualizar(context.Background(), AtualizarComandaPDVRequest{IDLoja: 0, IDComanda: 1, Finalizado: &finalizado})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when id_comanda is invalid", func(t *testing.T) {
		finalizado := true
		svc := NewComandaPDVService(lancamentoServiceStub{})

		_, err := svc.Atualizar(context.Background(), AtualizarComandaPDVRequest{IDLoja: 1, IDComanda: 0, Finalizado: &finalizado})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when finalizado is missing", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{})

		_, err := svc.Atualizar(context.Background(), AtualizarComandaPDVRequest{IDLoja: 1, IDComanda: 1})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("updates finalizado through lancamento service", func(t *testing.T) {
		finalizado := true
		svc := NewComandaPDVService(lancamentoServiceStub{updateFinalizadoFn: func(_ context.Context, req lancamento.UpdateFinalizadoRequest) (*models.LancamentoComanda, error) {
			if req.IDLoja != 1 || req.IDComanda != 2 || req.Finalizado == nil || !*req.Finalizado {
				t.Fatalf("unexpected update request: %+v", req)
			}
			return &models.LancamentoComanda{ID: 10, IDLoja: req.IDLoja, IDComanda: req.IDComanda, Finalizado: *req.Finalizado}, nil
		}})

		res, err := svc.Atualizar(context.Background(), AtualizarComandaPDVRequest{IDLoja: 1, IDComanda: 2, Finalizado: &finalizado})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.IDLoja != 1 || res.IDComanda != 2 || !res.Finalizado {
			t.Fatalf("unexpected response: %+v", res)
		}
	})

	t.Run("propagates lancamento update error", func(t *testing.T) {
		finalizado := true
		expectedErr := errors.New("update failed")
		svc := NewComandaPDVService(lancamentoServiceStub{updateFinalizadoFn: func(_ context.Context, _ lancamento.UpdateFinalizadoRequest) (*models.LancamentoComanda, error) {
			return nil, expectedErr
		}})

		_, err := svc.Atualizar(context.Background(), AtualizarComandaPDVRequest{IDLoja: 1, IDComanda: 2, Finalizado: &finalizado})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})
}
