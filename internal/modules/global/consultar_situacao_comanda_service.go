package global

import (
	"context"
	"fmt"
	"strings"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
)

type ConsultarSituacaoComandaService interface {
	Execute(ctx context.Context, req ConsultarSituacaoComandaRequest) (*ConsultarSituacaoComandaResponse, error)
}

type consultarSituacaoComandaService struct {
	lancamentoService lancamento.Service
	comandaService    comanda.Service
}

func NewConsultarSituacaoComandaService(lancamentoService lancamento.Service, comandaService comanda.Service) ConsultarSituacaoComandaService {
	return &consultarSituacaoComandaService{
		lancamentoService: lancamentoService,
		comandaService:    comandaService,
	}
}

func (s *consultarSituacaoComandaService) Execute(ctx context.Context, req ConsultarSituacaoComandaRequest) (*ConsultarSituacaoComandaResponse, error) {
	if req.IDLoja <= 0 {
		return nil, fmt.Errorf("%w: idLoja deve ser maior que zero", ErrInvalidRequest)
	}
	if strings.TrimSpace(req.NumeroIdentificacaoComanda) == "" {
		return nil, fmt.Errorf("%w: numeroIdentificacaoComanda e obrigatorio", ErrInvalidRequest)
	}

	comandas, err := s.comandaService.List(ctx, comanda.ListComandasRequest{
		IDLoja:              req.IDLoja,
		NumeroIdentificacao: req.NumeroIdentificacaoComanda,
	})
	if err != nil {
		return nil, err
	}
	if len(comandas) == 0 {
		return nil, nil
	}

	comanda := comandas[0]

	lancamentos, err := s.lancamentoService.List(ctx, lancamento.ListLancamentosRequest{
		IDLoja:    req.IDLoja,
		IDComanda: comanda.Comanda,
	})

	if err != nil {
		return nil, err
	}

	if len(lancamentos) == 0 {
		return &ConsultarSituacaoComandaResponse{
			IDLoja:                     req.IDLoja,
			Comanda:                    comanda.Comanda,
			NumeroIdentificacaoComanda: comanda.NumeroIdentificacao,
			Situacao:                   SituacaoComandaLiberada,
		}, nil
	}

	situacao := SituacaoComandaLiberada

	for _, lancamento := range lancamentos {
		if !lancamento.Finalizado {
			situacao = SituacaoComandaBloqueada
			break
		}
	}

	return &ConsultarSituacaoComandaResponse{
		IDLoja:                     req.IDLoja,
		Comanda:                    comanda.Comanda,
		NumeroIdentificacaoComanda: comanda.NumeroIdentificacao,
		Situacao:                   situacao,
	}, nil
}
