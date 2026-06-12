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

		_, err := svc.Consultar(context.Background(), GetLancamentoPDVRequest{NumeroComanda: 0, IDLoja: 1})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when loja is invalid", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{})

		_, err := svc.Consultar(context.Background(), GetLancamentoPDVRequest{NumeroComanda: 2, IDLoja: 0})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns not found when comanda does not exist", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{})

		res, err := svc.Consultar(context.Background(), GetLancamentoPDVRequest{
			NumeroComanda: 2,
			IDLoja:        1,
		})

		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
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

		res, err := svc.Consultar(context.Background(), GetLancamentoPDVRequest{NumeroComanda: 2, IDLoja: 1})
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

		_, err := svc.Consultar(context.Background(), GetLancamentoPDVRequest{NumeroComanda: 2, IDLoja: 1})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestComandaPDVServiceAtualizar(t *testing.T) {
	t.Run("returns invalid request when id_loja is invalid", func(t *testing.T) {
		finalizado := true
		svc := NewComandaPDVService(lancamentoServiceStub{})

		err := svc.Atualizar(context.Background(), UpdadeLancamentoPDVRequest{
			IDLoja:     0,
			IDComanda:  []int{1},
			Finalizado: &finalizado,
		})

		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when id_comanda is empty", func(t *testing.T) {
		finalizado := true
		svc := NewComandaPDVService(lancamentoServiceStub{})

		err := svc.Atualizar(context.Background(), UpdadeLancamentoPDVRequest{
			IDLoja:     1,
			IDComanda:  []int{},
			Finalizado: &finalizado,
		})

		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when finalizado is missing", func(t *testing.T) {
		svc := NewComandaPDVService(lancamentoServiceStub{})

		err := svc.Atualizar(context.Background(), UpdadeLancamentoPDVRequest{
			IDLoja:    1,
			IDComanda: []int{1},
		})

		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("propagates lancamento update error", func(t *testing.T) {
		finalizado := true
		expectedErr := errors.New("update failed")

		svc := NewComandaPDVService(lancamentoServiceStub{
			updateLancamentoByPDVFn: func(_ context.Context, _ lancamento.UpdateLancamentoByPDVRequest) error {
				return expectedErr
			},
		})

		err := svc.Atualizar(context.Background(), UpdadeLancamentoPDVRequest{
			IDLoja:     1,
			IDComanda:  []int{2, 3},
			Finalizado: &finalizado,
		})

		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("updates finalizado through lancamento service", func(t *testing.T) {
		finalizado := true

		var capturedReq lancamento.UpdateLancamentoByPDVRequest

		svc := NewComandaPDVService(lancamentoServiceStub{
			updateLancamentoByPDVFn: func(_ context.Context, req lancamento.UpdateLancamentoByPDVRequest) error {
				capturedReq = req
				return nil
			},
		})

		err := svc.Atualizar(context.Background(), UpdadeLancamentoPDVRequest{
			IDLoja:     1,
			IDComanda:  []int{2, 3},
			Finalizado: &finalizado,
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if capturedReq.IDLoja != 1 {
			t.Fatalf("expected id_loja=1, got %d", capturedReq.IDLoja)
		}

		if len(capturedReq.IDComanda) != 2 {
			t.Fatalf("expected 2 comandas, got %+v", capturedReq.IDComanda)
		}

		if capturedReq.IDComanda[0] != 2 || capturedReq.IDComanda[1] != 3 {
			t.Fatalf("unexpected comandas %+v", capturedReq.IDComanda)
		}

		if capturedReq.Finalizado == nil || !*capturedReq.Finalizado {
			t.Fatalf("expected finalizado=true")
		}
	})
}
