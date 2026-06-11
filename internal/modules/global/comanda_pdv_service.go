package global

import (
	"context"
	"fmt"

	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/models"
)

type ComandaPDVService interface {
	Consultar(ctx context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error)
	Atualizar(ctx context.Context, req AtualizarComandaPDVRequest) (*models.LancamentoComanda, error)
}

type comandaPDVService struct {
	lancamentoService lancamento.Service
}

func NewComandaPDVService(lancamentoService lancamento.Service) ComandaPDVService {
	return &comandaPDVService{lancamentoService: lancamentoService}
}

// Endpoint utilizado pelo PDV para consultar a comanda aberta.
func (s *comandaPDVService) Consultar(ctx context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error) {
	if req.NumeroComanda <= 0 {
		return nil, fmt.Errorf("%w: numeroComanda deve ser maior que zero", ErrInvalidRequest)
	}
	idLoja := req.IDLoja
	if idLoja <= 0 {
		idLoja = req.IDLojaAlias
	}
	if idLoja <= 0 {
		return nil, fmt.Errorf("%w: loja deve ser maior que zero", ErrInvalidRequest)
	}

	finalizado := false
	lancamentos, err := s.lancamentoService.List(ctx, lancamento.ListLancamentosRequest{
		IDLoja:     idLoja,
		IDComanda:  req.NumeroComanda,
		Finalizado: &finalizado,
	})
	if err != nil {
		return nil, err
	}
	if len(lancamentos) == 0 {
		return nil, nil
	}

	response := &ConsultarComandaPDVResponse{
		CodigoComanda:        lancamentos[0].IDComanda,
		TipoDocumentoCliente: 1,
		DocumentoCliente:     "",
		NomeCliente:          "",
		CodigoVendedor:       lancamentos[0].IDAtendente,
		ValorDescontoVenda:   0,
		ValorAcrescimoVenda:  0,
		Itens:                []ConsultarComandaPDVItemDTO{},
	}

	for _, lancamentoItem := range lancamentos {
		for _, item := range lancamentoItem.Itens {
			if item.Cancelado {
				continue
			}
			response.Itens = append(response.Itens, ConsultarComandaPDVItemDTO{
				CodigoBarras:   item.CodigoBarras,
				Quantidade:     item.Quantidade,
				PrecoVenda:     item.PrecoVenda,
				ValorDesconto:  0,
				ValorAcrescimo: 0,
			})
		}
	}

	return response, nil
}

// Endpoint utilizado pelo PDV para atualizar a situacao da comanda.
func (s *comandaPDVService) Atualizar(ctx context.Context, req AtualizarComandaPDVRequest) (*models.LancamentoComanda, error) {
	if req.IDLoja <= 0 {
		return nil, fmt.Errorf("%w: id_loja deve ser maior que zero", ErrInvalidRequest)
	}
	if req.IDComanda <= 0 {
		return nil, fmt.Errorf("%w: id_comanda deve ser maior que zero", ErrInvalidRequest)
	}
	if req.Finalizado == nil {
		return nil, fmt.Errorf("%w: finalizado e obrigatorio", ErrInvalidRequest)
	}

	result, err := s.lancamentoService.UpdateFinalizado(ctx, lancamento.UpdateFinalizadoRequest{
		IDLoja:     req.IDLoja,
		IDComanda:  req.IDComanda,
		Finalizado: req.Finalizado,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
