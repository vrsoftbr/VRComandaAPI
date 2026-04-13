package lancamento

import (
	"context"

	"gorm.io/gorm"

	"vrcomandaapi/internal/shared/models"
)

type Repository interface {
	Create(ctx context.Context, model *models.LancamentoComanda) error
	ExistsByLojaComanda(ctx context.Context, idLoja, idComanda int, finalizado bool) (bool, error)
	ExistsByLojaComandaExcludingID(ctx context.Context, id uint, idLoja, idComanda int, finalizado bool) (bool, error)
	FindByID(ctx context.Context, id uint) (*models.LancamentoComanda, error)
	Update(ctx context.Context, model *models.LancamentoComanda) error
	List(ctx context.Context, filter ListLancamentosFilter) ([]models.LancamentoComanda, error)
	ListItensByComanda(ctx context.Context, idComanda int) ([]ItemComandaRow, error)
	CreateItemsBatch(ctx context.Context, items []*models.LancamentoComandaItem) error
	SequenciaExistsInLancamento(ctx context.Context, idLancamento uint, sequencia int) (bool, error)
	FindItemByID(ctx context.Context, id uint) (*models.LancamentoComandaItem, error)
	UpdateItem(ctx context.Context, item *models.LancamentoComandaItem) error
}

type ItemComandaRow struct {
	models.LancamentoComandaItem
	IDComanda int `gorm:"column:id_comanda"`
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, model *models.LancamentoComanda) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository) ExistsByLojaComanda(ctx context.Context, idLoja, idComanda int, finalizado bool) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.LancamentoComanda{}).
		Where("id_loja = ? AND id_comanda = ? AND finalizado = ?", idLoja, idComanda, finalizado).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) ExistsByLojaComandaExcludingID(ctx context.Context, id uint, idLoja, idComanda int, finalizado bool) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.LancamentoComanda{}).
		Where("id <> ? AND id_loja = ? AND id_comanda = ? AND finalizado = ?", id, idLoja, idComanda, finalizado).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) FindByID(ctx context.Context, id uint) (*models.LancamentoComanda, error) {
	var model models.LancamentoComanda
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *repository) Update(ctx context.Context, model *models.LancamentoComanda) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *repository) CreateItemsBatch(ctx context.Context, items []*models.LancamentoComandaItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) SequenciaExistsInLancamento(ctx context.Context, idLancamento uint, sequencia int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.LancamentoComandaItem{}).
		Where("id_lancamentocomanda = ? AND sequencia = ?", idLancamento, sequencia).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) FindItemByID(ctx context.Context, id uint) (*models.LancamentoComandaItem, error) {
	var item models.LancamentoComandaItem
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) UpdateItem(ctx context.Context, item *models.LancamentoComandaItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *repository) ListItensByComanda(ctx context.Context, idComanda int) ([]ItemComandaRow, error) {
	var rows []ItemComandaRow
	err := r.db.WithContext(ctx).
		Table("lancamentocomandaitem AS i").
		Select("i.*, l.id_comanda").
		Joins("INNER JOIN lancamentocomanda AS l ON l.id = i.id_lancamentocomanda").
		Where("l.id_comanda = ?", idComanda).
		Order("i.id_lancamentocomanda, i.sequencia").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) List(ctx context.Context, filter ListLancamentosFilter) ([]models.LancamentoComanda, error) {
	query := r.db.WithContext(ctx).Model(&models.LancamentoComanda{})

	if filter.IDComanda != nil {
		query = query.Where("id_comanda = ?", *filter.IDComanda)
	}
	if filter.IDMesa != nil {
		query = query.Where("id_mesa = ?", *filter.IDMesa)
	}
	if filter.IDAtendente != nil {
		query = query.Where("id_atendente = ?", *filter.IDAtendente)
	}
	if filter.DataHora != nil {
		query = query.Where("data_hora = ?", *filter.DataHora)
	}
	if filter.Finalizado != nil {
		query = query.Where("finalizado = ?", *filter.Finalizado)
	}

	var result []models.LancamentoComanda
	if err := query.Order("id desc").Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
