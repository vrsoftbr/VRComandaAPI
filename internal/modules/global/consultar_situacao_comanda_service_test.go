package global

import (
	"context"
	"errors"
	"testing"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/models"
)

func TestConsultarSituacaoComandaServiceExecute(t *testing.T) {
	t.Run("returns invalid request when idLoja is invalid", func(t *testing.T) {
		svc := NewConsultarSituacaoComandaService(lancamentoServiceStub{}, comandaServiceStub{})

		_, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 0, NumeroIdentificacaoComanda: "100"})
		if err == nil || !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when numeroIdentificacaoComanda is empty", func(t *testing.T) {
		svc := NewConsultarSituacaoComandaService(lancamentoServiceStub{}, comandaServiceStub{})

		_, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "  "})
		if err == nil || !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns nil when comanda does not exist", func(t *testing.T) {
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return nil, nil
			}},
		)
		res, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "100"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})

	t.Run("returns liberada when comanda has no lancamentos", func(t *testing.T) {
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return nil, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)
		res, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.Situacao != SituacaoComandaLiberada {
			t.Fatalf("expected liberada, got %+v", res)
		}
	})

	t.Run("returns bloqueada when any lancamento is not finalized", func(t *testing.T) {
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{listFn: func(_ context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDComanda: 115, Finalizado: false}, {ID: 2, IDComanda: 115, Finalizado: true}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, req comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)
		res, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.Situacao != SituacaoComandaBloqueada {
			t.Fatalf("expected bloqueada, got %+v", res)
		}
	})

	t.Run("returns released when all lancamentos are finalized", func(t *testing.T) {
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDComanda: 115, Finalizado: true}, {ID: 2, IDComanda: 115, Finalizado: true}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)

		res, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.Situacao != SituacaoComandaLiberada {
			t.Fatalf("unexpected result: %+v", res)
		}
	})

	t.Run("returns liberada when all lancamentos are finalizados", func(t *testing.T) {
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDComanda: 115, Finalizado: true}, {ID: 2, IDComanda: 115, Finalizado: true}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)
		res, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.Situacao != SituacaoComandaLiberada {
			t.Fatalf("expected liberada, got %+v", res)
		}
	})
	t.Run("returns error when comanda service fails", func(t *testing.T) {
		expectedErr := errors.New("comanda service failed")
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return nil, expectedErr
			}},
		)
		_, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "100"})
		if err == nil || err.Error() != expectedErr.Error() {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("returns error when lancamento service fails", func(t *testing.T) {
		expectedErr := errors.New("lancamento service failed")
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return nil, expectedErr
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)
		_, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err == nil || err.Error() != expectedErr.Error() {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("propagates lancamento service error", func(t *testing.T) {
		expectedErr := errors.New("lancamento service failed")
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return nil, expectedErr
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "100"}}, nil
			}},
		)

		_, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "100"})
		if err == nil || err.Error() != expectedErr.Error() {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("propagates comanda service error", func(t *testing.T) {
		expectedErr := errors.New("comanda service failed")
		svc := NewConsultarSituacaoComandaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDComanda: 115, Finalizado: false}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return nil, expectedErr
			}},
		)

		_, err := svc.Execute(context.Background(), ConsultarSituacaoComandaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "100"})
		if err == nil || err.Error() != expectedErr.Error() {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})
}
