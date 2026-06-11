package global

import (
	"context"
	"fmt"
	"strings"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
)

type ConsultarComandaCatracaService interface {
	Execute(ctx context.Context, req ConsultarComandaCatracaRequest) (*ConsultarComandaCatracaResponse, error)
}

type consultarComandaCatracaService struct {
	lancamentoService lancamento.Service
	comandaService    comanda.Service
}

func NewConsultarComandaCatracaService(lancamentoService lancamento.Service, comandaService comanda.Service) ConsultarComandaCatracaService {
	return &consultarComandaCatracaService{
		lancamentoService: lancamentoService,
		comandaService:    comandaService,
	}
}

// Endpoint utilizado pela aplicação da Catraca
func (s *consultarComandaCatracaService) Execute(ctx context.Context, req ConsultarComandaCatracaRequest) (*ConsultarComandaCatracaResponse, error) {
	if req.IDLoja <= 0 {
		return nil, fmt.Errorf("%w: idLoja deve ser maior que zero", ErrInvalidRequest)
	}
	if strings.TrimSpace(req.NumeroIdentificacaoComanda) == "" {
		return nil, fmt.Errorf("%w: numeroIdentificacaoComanda e obrigatorio", ErrInvalidRequest)
	}

	lancamentos, err := s.lancamentoService.List(ctx, lancamento.ListLancamentosRequest{
		IDLoja: req.IDLoja,
	})
	if err != nil {
		return nil, err
	}
	if len(lancamentos) == 0 {
		return nil, nil
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

	detalheComanda := comandas[0]

	lancamentosComanda := make([]bool, 0, len(lancamentos))
	for _, item := range lancamentos {
		if item.IDComanda == detalheComanda.Comanda {
			lancamentosComanda = append(lancamentosComanda, item.Finalizado)
		}
	}
	if len(lancamentosComanda) == 0 {
		return nil, nil
	}

	situacao := SituacaoComandaLiberada
	for _, finalizado := range lancamentosComanda {
		if !finalizado {
			situacao = SituacaoComandaBloqueada
			break
		}
	}

	return &ConsultarComandaCatracaResponse{
		IDLoja:                     req.IDLoja,
		Comanda:                    detalheComanda.Comanda,
		NumeroIdentificacaoComanda: detalheComanda.NumeroIdentificacao,
		Situacao:                   situacao,
	}, nil
}
