package lancamento

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"vrcomandaapi/internal/shared/models"
)

var ErrInvalidFilter = errors.New("filtro invalido")
var ErrValidation = errors.New("erro de validacao")
var ErrDuplicateLancamento = errors.New("ja existe lancamento com id_loja e id_comanda informados")
var ErrNotFound = errors.New("lancamento nao encontrado")

var ErrDuplicateSequencia = errors.New("sequencia ja existe nesse lancamento")
var ErrComandaFinalizada = errors.New("comanda ja esta finalizada")

type Service interface {
	Create(ctx context.Context, req CreateLancamentoRequest) (*models.LancamentoComanda, error)
	Update(ctx context.Context, id uint, req CreateLancamentoRequest) (*models.LancamentoComanda, error)
	List(ctx context.Context, req ListLancamentosRequest) ([]models.LancamentoComanda, error)
	ListItens(ctx context.Context, req ListItensRequest) ([]ItemComandaResponse, error)
	CreateItems(ctx context.Context, req CreateItemsRequest) ([]*models.LancamentoComandaItem, error)
	UpdateItem(ctx context.Context, id uint, req UpdateItemRequest) (*models.LancamentoComandaItem, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateLancamentoRequest) (*models.LancamentoComanda, error) {
	dataHora, err := validateCreateLancamentoRequest(req)
	if err != nil {
		return nil, err
	}

	if !*req.Finalizado {
		exists, err := s.repo.ExistsByLojaComanda(ctx, req.IDLoja, req.IDComanda, false)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrDuplicateLancamento
		}
	}

	model := &models.LancamentoComanda{}
	applyLancamentoRequest(model, req, dataHora)

	if err := s.repo.Create(ctx, model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *service) List(ctx context.Context, req ListLancamentosRequest) ([]models.LancamentoComanda, error) {
	filter := ListLancamentosFilter{}

	if req.IDLoja != 0 {
		v := req.IDLoja
		filter.IDLoja = &v
	}

	if req.IDComanda != 0 {
		v := req.IDComanda
		filter.IDComanda = &v
	}

	if req.IDMesa != 0 {
		v := req.IDMesa
		filter.IDMesa = &v
	}

	if req.IDAtendente != 0 {
		v := req.IDAtendente
		filter.IDAtendente = &v
	}

	if req.Finalizado != nil {
		filter.Finalizado = req.Finalizado
	}

	if req.DataHora != "" {
		parsed, err := parseDataHora(req.DataHora)
		if err != nil {
			return nil, fmt.Errorf("%w: dataHora invalida", ErrInvalidFilter)
		}
		filter.DataHora = &parsed
	}

	return s.repo.List(ctx, filter)
}

func (s *service) Update(ctx context.Context, id uint, req CreateLancamentoRequest) (*models.LancamentoComanda, error) {
	dataHora, err := validateCreateLancamentoRequest(req)
	if err != nil {
		return nil, err
	}

	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if model.Finalizado && !*req.Finalizado {
		return nil, fmt.Errorf("%w: nao e permitido alterar finalizado de true para false", ErrValidation)
	}

	if !*req.Finalizado {
		exists, err := s.repo.ExistsByLojaComandaExcludingID(ctx, id, req.IDLoja, req.IDComanda, false)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrDuplicateLancamento
		}
	}

	applyLancamentoRequest(model, req, dataHora)

	if err := s.repo.Update(ctx, model); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *service) CreateItems(ctx context.Context, req CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
	seen := make(map[string]struct{}, len(req.Itens))

	items := make([]*models.LancamentoComandaItem, 0, len(req.Itens))
	for i, r := range req.Itens {
		if r.Quantidade <= 0 {
			return nil, fmt.Errorf("%w: item %d: quantidade deve ser maior que zero", ErrValidation, i+1)
		}
		if r.PrecoVenda < 0 {
			return nil, fmt.Errorf("%w: item %d: precovenda nao pode ser negativo", ErrValidation, i+1)
		}

		lanc, err := s.repo.FindByID(ctx, r.IDLancamentoComanda)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("%w: item %d: id_lancamentocomanda nao encontrado", ErrNotFound, i+1)
			}
			return nil, err
		}
		if lanc.Finalizado {
			return nil, fmt.Errorf("%w: item %d", ErrComandaFinalizada, i+1)
		}

		duplicado, err := s.repo.SequenciaExistsInLancamento(ctx, r.IDLancamentoComanda, r.Sequencia)
		if err != nil {
			return nil, err
		}
		if duplicado {
			return nil, fmt.Errorf("%w: item %d: sequencia %d", ErrDuplicateSequencia, i+1, r.Sequencia)
		}

		key := fmt.Sprintf("%d:%d", r.IDLancamentoComanda, r.Sequencia)
		if _, exists := seen[key]; exists {
			return nil, fmt.Errorf("%w: item %d: sequencia %d duplicada no mesmo lote", ErrDuplicateSequencia, i+1, r.Sequencia)
		}
		seen[key] = struct{}{}

		items = append(items, &models.LancamentoComandaItem{
			IDLancamentoComanda: r.IDLancamentoComanda,
			Sequencia:           r.Sequencia,
			IDProduto:           r.IDProduto,
			CodigoBarras:        r.CodigoBarras,
			Quantidade:          r.Quantidade,
			PrecoVenda:          r.PrecoVenda,
			IDAtendente:         r.IDAtendente,
			IDSituacao:          r.IDSituacao,
		})
	}

	if err := s.repo.CreateItemsBatch(ctx, items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *service) UpdateItem(ctx context.Context, id uint, req UpdateItemRequest) (*models.LancamentoComandaItem, error) {
	if req.Quantidade == nil && req.Cancelado == nil {
		return nil, fmt.Errorf("%w: informe quantidade ou cancelado", ErrValidation)
	}
	if req.Quantidade != nil && *req.Quantidade <= 0 {
		return nil, fmt.Errorf("%w: quantidade deve ser maior que zero", ErrValidation)
	}

	item, err := s.repo.FindItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	lanc, err := s.repo.FindByID(ctx, item.IDLancamentoComanda)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if lanc.Finalizado {
		return nil, ErrComandaFinalizada
	}

	if req.Quantidade != nil {
		item.Quantidade = *req.Quantidade
	}

	if req.Cancelado != nil {
		item.Cancelado = *req.Cancelado
		if *req.Cancelado {
			now := time.Now()
			item.DataHoraCancelamento = &now
		} else {
			item.DataHoraCancelamento = nil
		}
	}

	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) ListItens(ctx context.Context, req ListItensRequest) ([]ItemComandaResponse, error) {
	rows, err := s.repo.ListItensByComanda(ctx, req.IDComanda)
	if err != nil {
		return nil, err
	}

	response := make([]ItemComandaResponse, 0, len(rows))
	for _, r := range rows {
		response = append(response, ItemComandaResponse{
			ID:                   r.ID,
			IDLancamentoComanda:  r.IDLancamentoComanda,
			IDComanda:            r.IDComanda,
			Sequencia:            r.Sequencia,
			IDProduto:            r.IDProduto,
			CodigoBarras:         r.CodigoBarras,
			Quantidade:           r.Quantidade,
			PrecoVenda:           r.PrecoVenda,
			Cancelado:            r.Cancelado,
			DataHoraCancelamento: r.DataHoraCancelamento,
			IDAtendente:          r.IDAtendente,
			IDSituacao:           r.IDSituacao,
		})
	}

	return response, nil
}

func parseDataHora(input string) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, input); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("dataHora invalida")
}

func validateCreateLancamentoRequest(req CreateLancamentoRequest) (time.Time, error) {
	if req.IDLoja <= 0 {
		return time.Time{}, fmt.Errorf("%w: id_loja e obrigatorio", ErrValidation)
	}
	if req.IDComanda <= 0 {
		return time.Time{}, fmt.Errorf("%w: id_comanda e obrigatorio", ErrValidation)
	}
	if req.IDAtendente <= 0 {
		return time.Time{}, fmt.Errorf("%w: id_atendente e obrigatorio", ErrValidation)
	}
	if req.Finalizado == nil {
		return time.Time{}, fmt.Errorf("%w: finalizado e obrigatorio", ErrValidation)
	}
	if strings.TrimSpace(req.DataHora) == "" {
		return time.Time{}, fmt.Errorf("%w: dataHora e obrigatorio", ErrValidation)
	}

	dataHora, err := parseDataHora(req.DataHora)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: dataHora invalida", ErrValidation)
	}

	return dataHora, nil
}

func applyLancamentoRequest(model *models.LancamentoComanda, req CreateLancamentoRequest, dataHora time.Time) {
	model.IDLoja = req.IDLoja
	model.IDComanda = req.IDComanda
	model.IDMesa = req.IDMesa
	model.IDAtendente = req.IDAtendente
	model.DataHora = dataHora
	model.Observacao = strings.TrimSpace(req.Observacao)
	model.Finalizado = *req.Finalizado
}
