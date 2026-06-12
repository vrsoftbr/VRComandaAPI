package global

import (
	"context"
	"fmt"

	"vrcomandaapi/internal/modules/lancamento"
)

type ComandaPDVService interface {
	Consultar(ctx context.Context, req GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error)
	Atualizar(ctx context.Context, req UpdadeLancamentoPDVRequest) error
}

type comandaPDVService struct {
	lancamentoService lancamento.Service
}

func NewComandaPDVService(lancamentoService lancamento.Service) ComandaPDVService {
	return &comandaPDVService{lancamentoService: lancamentoService}
}

func (s *comandaPDVService) Consultar(ctx context.Context, req GetLancamentoPDVRequest) (*GetLancamentoItemPDVResponse, error) {
	if req.NumeroComanda <= 0 {
		return nil, fmt.Errorf("%w: numeroComanda deve ser maior que zero", ErrInvalidRequest)
	}

	if req.IDLoja <= 0 {
		return nil, fmt.Errorf("%w: idLoja deve ser maior que zero", ErrInvalidRequest)
	}

	finalizado := false
	lancamentos, err := s.lancamentoService.List(ctx, lancamento.ListLancamentosRequest{
		IDLoja:     req.IDLoja,
		IDComanda:  req.NumeroComanda,
		Finalizado: &finalizado,
	})
	if err != nil {
		return nil, err
	}
	if len(lancamentos) == 0 {
		return nil, fmt.Errorf("%w: comanda nao encontrada", ErrNotFound)
	}

	response := &GetLancamentoItemPDVResponse{
		CodigoComanda:        lancamentos[0].IDComanda,
		TipoDocumentoCliente: 1,
		DocumentoCliente:     "",
		NomeCliente:          "",
		CodigoVendedor:       lancamentos[0].IDAtendente,
		ValorDescontoVenda:   0,
		ValorAcrescimoVenda:  0,
		Itens:                []GetLancamentoPDVItemDTO{},
	}

	for _, lancamentoItem := range lancamentos {
		for _, item := range lancamentoItem.Itens {
			if item.Cancelado {
				continue
			}
			response.Itens = append(response.Itens, GetLancamentoPDVItemDTO{
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
func (s *comandaPDVService) Atualizar(ctx context.Context, req UpdadeLancamentoPDVRequest) error {
	if req.IDLoja <= 0 {
		return fmt.Errorf("%w: id_loja deve ser maior que zero", ErrInvalidRequest)
	}
	if len(req.IDComanda) == 0 {
		return fmt.Errorf("%w: id_comanda deve conter pelo menos um id de comanda", ErrInvalidRequest)
	}
	if req.Finalizado == nil {
		return fmt.Errorf("%w: finalizado e obrigatorio", ErrInvalidRequest)
	}

	if err := s.lancamentoService.UpdateLancamentoByPDV(ctx, lancamento.UpdateLancamentoByPDVRequest{
		IDLoja:     req.IDLoja,
		IDComanda:  req.IDComanda,
		Finalizado: req.Finalizado,
	}); err != nil {
		return err
	}

	return nil
}
