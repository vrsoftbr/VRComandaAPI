package global

import (
	"context"
	"fmt"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/models"
)

type ComandaPDVService interface {
	Consultar(ctx context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error)
	Atualizar(ctx context.Context, req AtualizarComandaPDVRequest) (*models.LancamentoComanda, error)
}

type comandaPDVService struct {
	lancamentoService lancamento.Service
	comandaService    comanda.Service
}

func NewComandaPDVService(lancamentoService lancamento.Service, comandaService comanda.Service) ComandaPDVService {
	return &comandaPDVService{lancamentoService: lancamentoService, comandaService: comandaService}
}

func (s *comandaPDVService) Consultar(ctx context.Context, req ConsultarComandaPDVRequest) (*ConsultarComandaPDVResponse, error) {
	if req.NumeroComanda == "" {
		return nil, fmt.Errorf("%w: numeroComanda deve ser maior que zero", ErrInvalidRequest)
	}
	idLoja := req.IDLoja
	if idLoja <= 0 {
		idLoja = req.IDLojaAlias
	}
	if idLoja <= 0 {
		return nil, fmt.Errorf("%w: loja deve ser maior que zero", ErrInvalidRequest)
	}

	comandas, err := s.comandaService.List(ctx, comanda.ListComandasRequest{
		IDLoja:              idLoja,
		NumeroIdentificacao: req.NumeroComanda,
	})

	if err != nil {
		return nil, err
	}
	if len(comandas) == 0 {
		return nil, nil
	}

	comanda := comandas[0]

	finalizado := false
	lancamentos, err := s.lancamentoService.List(ctx, lancamento.ListLancamentosRequest{
		IDLoja:     idLoja,
		IDComanda:  comanda.Comanda,
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
	if len(req.IDComanda) == 0 {
		return nil, fmt.Errorf("%w: id_comanda deve ser maior que zero", ErrInvalidRequest)
	}
	if req.Finalizado == nil {
		return nil, fmt.Errorf("%w: finalizado e obrigatorio", ErrInvalidRequest)
	}

	if err := s.lancamentoService.UpdateFinalizado(ctx, lancamento.UpdateFinalizadoRequest{
		IDLoja:     req.IDLoja,
		IDComanda:  req.IDComanda,
		Finalizado: req.Finalizado,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
