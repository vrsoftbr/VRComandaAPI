package global

import (
	"context"
	"sort"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/shared/models"
)

type LancamentosDetalhesService interface {
	Execute(ctx context.Context, req ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error)
}

type lancamentosDetalhesService struct {
	lancamentoService lancamento.Service
	comandaService    comanda.Service
	mesaService       mesa.Service
}

func NewLancamentosDetalhesService(lancamentoService lancamento.Service, comandaService comanda.Service, mesaService mesa.Service) LancamentosDetalhesService {
	return &lancamentosDetalhesService{
		lancamentoService: lancamentoService,
		comandaService:    comandaService,
		mesaService:       mesaService,
	}
}

func (s *lancamentosDetalhesService) Execute(ctx context.Context, req ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error) {
	finalizado := false
	if req.Finalizado != nil {
		finalizado = *req.Finalizado
	}

	lancamentos, err := s.lancamentoService.List(ctx, lancamento.ListLancamentosRequest{
		IDLoja:     req.IDLoja,
		Finalizado: &finalizado,
	})
	if err != nil {
		return nil, err
	}
	if len(lancamentos) == 0 {
		return []LancamentoDetalhesDTO{}, nil
	}

	comandas, err := s.comandaService.List(ctx, comanda.ListComandasRequest{
		IDLoja:   req.IDLoja,
		Comandas: uniqueSortedComandaIDs(lancamentos),
	})
	if err != nil {
		return nil, err
	}

	mesas, err := s.mesaService.List(ctx, mesa.ListMesasRequest{
		IDLoja: req.IDLoja,
		Mesas:  uniqueSortedMesaIDs(lancamentos),
	})
	if err != nil {
		return nil, err
	}

	comandasByID := buildComandasByID(comandas)
	mesasByID := buildMesasByID(mesas)

	result := make([]LancamentoDetalhesDTO, 0, len(lancamentos))
	for _, item := range lancamentos {
		row := LancamentoDetalhesDTO{
			IDLancamento: item.ID,
			IDLoja:       item.IDLoja,
			IDComanda:    item.IDComanda,
			IDMesa:       item.IDMesa,
			DataHora:     item.DataHora.Format("2006-01-02T15:04:05Z07:00"),
			Observacao:   item.Observacao,
			IDAtendente:  item.IDAtendente,
			Finalizado:   item.Finalizado,
			Comanda:      comandasByID[item.IDComanda],
			Itens:        cloneItens(item.Itens),
		}
		if item.IDMesa != nil {
			row.Mesa = mesasByID[*item.IDMesa]
		}
		result = append(result, row)
	}

	return result, nil
}

func uniqueSortedComandaIDs(lancamentos []models.LancamentoComanda) []int {
	seen := make(map[int]struct{}, len(lancamentos))
	ids := make([]int, 0, len(lancamentos))
	for _, item := range lancamentos {
		if _, ok := seen[item.IDComanda]; ok {
			continue
		}
		seen[item.IDComanda] = struct{}{}
		ids = append(ids, item.IDComanda)
	}
	sort.Ints(ids)
	return ids
}

func uniqueSortedMesaIDs(lancamentos []models.LancamentoComanda) []int {
	seen := make(map[int]struct{}, len(lancamentos))
	ids := make([]int, 0, len(lancamentos))
	for _, item := range lancamentos {
		if item.IDMesa == nil {
			continue
		}
		if _, ok := seen[*item.IDMesa]; ok {
			continue
		}
		seen[*item.IDMesa] = struct{}{}
		ids = append(ids, *item.IDMesa)
	}
	sort.Ints(ids)
	return ids
}

func buildComandasByID(items []comanda.ComandaResponse) map[int]*comanda.ComandaResponse {
	result := make(map[int]*comanda.ComandaResponse, len(items))
	for _, item := range items {
		copyItem := item
		result[item.Comanda] = &copyItem
	}
	return result
}

func buildMesasByID(items []mesa.MesaResponse) map[int]*mesa.MesaResponse {
	result := make(map[int]*mesa.MesaResponse, len(items))
	for _, item := range items {
		copyItem := item
		result[item.Mesa] = &copyItem
	}
	return result
}

func cloneItens(itens []models.LancamentoComandaItem) []models.LancamentoComandaItem {
	if len(itens) == 0 {
		return []models.LancamentoComandaItem{}
	}

	cloned := make([]models.LancamentoComandaItem, len(itens))
	copy(cloned, itens)
	return cloned
}
