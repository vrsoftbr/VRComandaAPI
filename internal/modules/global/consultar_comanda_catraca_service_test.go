package global

import (
	"context"
	"errors"
	"testing"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/models"
)

func TestConsultarComandaCatracaServiceExecute(t *testing.T) {
	t.Run("returns invalid request when idLoja is invalid", func(t *testing.T) {
		svc := NewConsultarComandaCatracaService(lancamentoServiceStub{}, comandaServiceStub{})

		_, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 0, NumeroIdentificacaoComanda: "100"})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns invalid request when numeroIdentificacaoComanda is empty", func(t *testing.T) {
		svc := NewConsultarComandaCatracaService(lancamentoServiceStub{}, comandaServiceStub{})

		_, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "  "})
		if !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected ErrInvalidRequest, got %v", err)
		}
	})

	t.Run("returns nil when comanda does not exist in mongo", func(t *testing.T) {
		svc := NewConsultarComandaCatracaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDLoja: 1, IDComanda: 115, Finalizado: false}}, nil
			}},
			comandaServiceStub{},
		)

		res, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "100"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})

	t.Run("returns nil when comanda has no lancamentos in sqlite", func(t *testing.T) {
		svc := NewConsultarComandaCatracaService(
			lancamentoServiceStub{},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)

		res, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})

	t.Run("returns blocked when any lancamento is not finalized", func(t *testing.T) {
		svc := NewConsultarComandaCatracaService(
			lancamentoServiceStub{listFn: func(_ context.Context, req lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				if req.IDLoja != 1 || req.IDComanda != 0 {
					t.Fatalf("unexpected request: %+v", req)
				}
				return []models.LancamentoComanda{{ID: 1, IDComanda: 114, Finalizado: false}, {ID: 2, IDComanda: 115, Finalizado: true}, {ID: 3, IDComanda: 115, Finalizado: false}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, req comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				if req.IDLoja != 1 || req.NumeroIdentificacao != "1000000001159" {
					t.Fatalf("unexpected request: %+v", req)
				}
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)

		res, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.Situacao != SituacaoComandaBloqueada {
			t.Fatalf("unexpected result: %+v", res)
		}
	})

	t.Run("returns released when all lancamentos are finalized", func(t *testing.T) {
		svc := NewConsultarComandaCatracaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDComanda: 115, Finalizado: true}, {ID: 2, IDComanda: 115, Finalizado: true}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)

		res, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.Situacao != SituacaoComandaLiberada {
			t.Fatalf("unexpected result: %+v", res)
		}
	})

	t.Run("returns nil when sqlite has records but not for searched comanda", func(t *testing.T) {
		svc := NewConsultarComandaCatracaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDComanda: 999, Finalizado: false}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return []comanda.ComandaResponse{{IDLoja: 1, Comanda: 115, NumeroIdentificacao: "1000000001159"}}, nil
			}},
		)

		res, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "1000000001159"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})

	t.Run("propagates lancamento service error", func(t *testing.T) {
		expectedErr := errors.New("lancamento service failed")
		svc := NewConsultarComandaCatracaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return nil, expectedErr
			}},
			comandaServiceStub{},
		)

		_, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "100"})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("propagates comanda service error", func(t *testing.T) {
		expectedErr := errors.New("comanda service failed")
		svc := NewConsultarComandaCatracaService(
			lancamentoServiceStub{listFn: func(_ context.Context, _ lancamento.ListLancamentosRequest) ([]models.LancamentoComanda, error) {
				return []models.LancamentoComanda{{ID: 1, IDComanda: 115, Finalizado: false}}, nil
			}},
			comandaServiceStub{listFn: func(_ context.Context, _ comanda.ListComandasRequest) ([]comanda.ComandaResponse, error) {
				return nil, expectedErr
			}},
		)

		_, err := svc.Execute(context.Background(), ConsultarComandaCatracaRequest{IDLoja: 1, NumeroIdentificacaoComanda: "100"})
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
	})
}
