package global

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"vrcomandaapi/internal/modules/atendente"
	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/modules/produto"
	"vrcomandaapi/internal/shared/models"
)

type LancamentosDetalhesService interface {
	Execute(ctx context.Context, req ListLancamentosDetalhesRequest) ([]LancamentoDetalhesDTO, error)
}

type lancamentosDetalhesService struct {
	lancamentoService lancamento.Service
	atendenteService  atendente.Service
	comandaService    comanda.Service
	mesaService       mesa.Service
	produtoService    produto.Service
}

func NewLancamentosDetalhesService(lancamentoService lancamento.Service, atendenteService atendente.Service, comandaService comanda.Service, mesaService mesa.Service, produtoService produto.Service) LancamentosDetalhesService {
	return &lancamentosDetalhesService{
		lancamentoService: lancamentoService,
		atendenteService:  atendenteService,
		comandaService:    comandaService,
		mesaService:       mesaService,
		produtoService:    produtoService,
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

	atendentesByKey, err := s.listAtendentesByLancamento(ctx, lancamentos)
	if err != nil {
		return nil, err
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
	produtosByBarcode, err := s.listProdutosByBarcode(ctx, req.IDLoja, lancamentos)
	if err != nil {
		return nil, err
	}

	result := make([]LancamentoDetalhesDTO, 0, len(lancamentos))
	for _, item := range lancamentos {

		row := LancamentoDetalhesDTO{
			IDLancamento: item.ID,
			IDLoja:       item.IDLoja,
			IDComanda:    item.IDComanda,
			IDMesa:       item.IDMesa,
			DataHora:     item.DataHora.Format(time.RFC3339),
			Observacao:   item.Observacao,
			IDAtendente:  item.IDAtendente,
			Finalizado:   item.Finalizado,
			Atendente:    atendentesByKey[buildAtendenteKey(item.IDLoja, item.IDAtendente)],
			Comanda:      comandasByID[item.IDComanda],
			Itens:        cloneItens(item.Itens, produtosByBarcode),
		}
		if item.IDMesa != nil {
			row.Mesa = mesasByID[*item.IDMesa]
		}
		result = append(result, row)
	}

	return result, nil
}

func (s *lancamentosDetalhesService) listAtendentesByLancamento(ctx context.Context, lancamentos []models.LancamentoComanda) (map[string]*atendente.AtendenteResponse, error) {
	keys := uniqueSortedAtendenteKeys(lancamentos)
	if len(keys) == 0 {
		return map[string]*atendente.AtendenteResponse{}, nil
	}

	atendentes := make([]atendente.AtendenteResponse, 0, len(keys))
	for _, key := range keys {
		response, err := s.atendenteService.List(ctx, atendente.ListAtendentesRequest{
			IDLoja:      key.idLoja,
			IDAtendente: strconv.Itoa(key.idAtendente),
		})
		if err != nil {
			return nil, err
		}
		if len(response) == 0 {
			continue
		}

		atendentes = append(atendentes, response[0])
	}

	return buildAtendentesByKey(atendentes), nil
}

func (s *lancamentosDetalhesService) listProdutosByBarcode(ctx context.Context, idLoja int, lancamentos []models.LancamentoComanda) (map[string]*produto.ProdutoResponse, error) {
	barras := uniqueSortedItemBarcodes(lancamentos)
	if len(barras) == 0 {
		return map[string]*produto.ProdutoResponse{}, nil
	}

	produtos := make([]produto.ProdutoResponse, 0, len(barras))
	for _, codigoBarras := range barras {
		responseAny, err := s.produtoService.List(ctx, produto.ListProdutosRequest{IDLoja: idLoja, CodigoBarras: codigoBarras})
		if err != nil {
			return nil, err
		}

		var response []produto.ProdutoResponse
		switch value := responseAny.(type) {
		case produto.ProdutosPaginatedResponse:
			response = value.Items
		case *produto.ProdutosPaginatedResponse:
			if value != nil {
				response = value.Items
			}
		case []produto.ProdutoResponse:
			response = value
		default:
			return nil, fmt.Errorf("unexpected produto response type: %T", responseAny)
		}

		if len(response) == 0 {
			continue
		}

		produtos = append(produtos, response[0])
	}

	return buildProdutosByBarcode(produtos), nil
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

type atendenteKey struct {
	idLoja      int
	idAtendente int
}

func uniqueSortedAtendenteKeys(lancamentos []models.LancamentoComanda) []atendenteKey {
	seen := make(map[atendenteKey]struct{}, len(lancamentos))
	keys := make([]atendenteKey, 0, len(lancamentos))
	for _, item := range lancamentos {
		key := atendenteKey{idLoja: item.IDLoja, idAtendente: item.IDAtendente}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].idLoja == keys[j].idLoja {
			return keys[i].idAtendente < keys[j].idAtendente
		}
		return keys[i].idLoja < keys[j].idLoja
	})
	return keys
}

func buildAtendenteKey(idLoja int, idAtendente int) string {
	return strconv.Itoa(idLoja) + ":" + strconv.Itoa(idAtendente)
}

func buildAtendentesByKey(items []atendente.AtendenteResponse) map[string]*atendente.AtendenteResponse {
	result := make(map[string]*atendente.AtendenteResponse, len(items))
	for _, item := range items {
		idAtendente, err := strconv.Atoi(item.IDAtendente)
		if err != nil {
			continue
		}
		copyItem := item
		result[buildAtendenteKey(item.IDLoja, idAtendente)] = &copyItem
	}
	return result
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

func buildProdutosByBarcode(items []produto.ProdutoResponse) map[string]*produto.ProdutoResponse {
	result := make(map[string]*produto.ProdutoResponse, len(items))

	for _, item := range items {
		if strings.TrimSpace(item.CodigoBarras) == "" {
			continue
		}
		copyItem := item
		result[item.CodigoBarras] = &copyItem
	}

	return result
}

func uniqueSortedItemBarcodes(lancamentos []models.LancamentoComanda) []string {
	seen := make(map[string]struct{})
	barras := make([]string, 0)

	for _, lancamento := range lancamentos {
		for _, item := range lancamento.Itens {
			codigoBarras := strings.TrimSpace(item.CodigoBarras)
			if codigoBarras == "" {
				continue
			}
			if _, ok := seen[codigoBarras]; ok {
				continue
			}
			seen[codigoBarras] = struct{}{}
			barras = append(barras, codigoBarras)
		}
	}

	sort.Strings(barras)
	return barras
}

func cloneItens(itens []models.LancamentoComandaItem, produtosByBarcode map[string]*produto.ProdutoResponse) []LancamentoDetalhesItemDTO {
	if len(itens) == 0 {
		return []LancamentoDetalhesItemDTO{}
	}

	cloned := make([]LancamentoDetalhesItemDTO, 0, len(itens))
	for _, item := range itens {
		row := LancamentoDetalhesItemDTO{LancamentoComandaItem: item}
		if produtoItem, ok := produtosByBarcode[item.CodigoBarras]; ok {
			row.DescricaoProduto = produtoItem.DescricaoCompleta
			if row.DescricaoProduto == "" {
				row.DescricaoProduto = produtoItem.DescricaoCupom
			}
		}

		cloned = append(cloned, row)
	}

	return cloned
}
