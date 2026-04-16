package lancamento

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"vrcomandaapi/internal/shared/models"
)

type repositoryStub struct {
	createFn                         func(ctx context.Context, model *models.LancamentoComanda) error
	existsByLojaComandaFn            func(ctx context.Context, idLoja, idComanda int, finalizado bool) (bool, error)
	existsByLojaComandaExcludingIDFn func(ctx context.Context, id uint, idLoja, idComanda int, finalizado bool) (bool, error)
	findByIDFn                       func(ctx context.Context, id uint) (*models.LancamentoComanda, error)
	updateFn                         func(ctx context.Context, model *models.LancamentoComanda) error
	listFn                           func(ctx context.Context, filter ListLancamentosFilter) ([]models.LancamentoComanda, error)
	listItensByComandaFn             func(ctx context.Context, idComanda int) ([]ItemComandaRow, error)
	createItemsBatchFn               func(ctx context.Context, items []*models.LancamentoComandaItem) error
	sequenciaExistsInLancamentoFn    func(ctx context.Context, idLancamento uint, sequencia int) (bool, error)
	findItemByIDFn                   func(ctx context.Context, id uint) (*models.LancamentoComandaItem, error)
	updateItemFn                     func(ctx context.Context, item *models.LancamentoComandaItem) error
}

func (s repositoryStub) Create(ctx context.Context, model *models.LancamentoComanda) error {
	if s.createFn != nil {
		return s.createFn(ctx, model)
	}
	return nil
}

func (s repositoryStub) ExistsByLojaComanda(ctx context.Context, idLoja, idComanda int, finalizado bool) (bool, error) {
	if s.existsByLojaComandaFn != nil {
		return s.existsByLojaComandaFn(ctx, idLoja, idComanda, finalizado)
	}
	return false, nil
}

func (s repositoryStub) ExistsByLojaComandaExcludingID(ctx context.Context, id uint, idLoja, idComanda int, finalizado bool) (bool, error) {
	if s.existsByLojaComandaExcludingIDFn != nil {
		return s.existsByLojaComandaExcludingIDFn(ctx, id, idLoja, idComanda, finalizado)
	}
	return false, nil
}

func (s repositoryStub) FindByID(ctx context.Context, id uint) (*models.LancamentoComanda, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s repositoryStub) Update(ctx context.Context, model *models.LancamentoComanda) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, model)
	}
	return nil
}

func (s repositoryStub) List(ctx context.Context, filter ListLancamentosFilter) ([]models.LancamentoComanda, error) {
	if s.listFn != nil {
		return s.listFn(ctx, filter)
	}
	return []models.LancamentoComanda{}, nil
}

func (s repositoryStub) ListItensByComanda(ctx context.Context, idComanda int) ([]ItemComandaRow, error) {
	if s.listItensByComandaFn != nil {
		return s.listItensByComandaFn(ctx, idComanda)
	}
	return []ItemComandaRow{}, nil
}

func (s repositoryStub) CreateItemsBatch(ctx context.Context, items []*models.LancamentoComandaItem) error {
	if s.createItemsBatchFn != nil {
		return s.createItemsBatchFn(ctx, items)
	}
	return nil
}

func (s repositoryStub) SequenciaExistsInLancamento(ctx context.Context, idLancamento uint, sequencia int) (bool, error) {
	if s.sequenciaExistsInLancamentoFn != nil {
		return s.sequenciaExistsInLancamentoFn(ctx, idLancamento, sequencia)
	}
	return false, nil
}

func (s repositoryStub) FindItemByID(ctx context.Context, id uint) (*models.LancamentoComandaItem, error) {
	if s.findItemByIDFn != nil {
		return s.findItemByIDFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s repositoryStub) UpdateItem(ctx context.Context, item *models.LancamentoComandaItem) error {
	if s.updateItemFn != nil {
		return s.updateItemFn(ctx, item)
	}
	return nil
}

func boolPtr(v bool) *bool        { return &v }
func floatPtr(v float64) *float64 { return &v }

func validCreateRequest() CreateLancamentoRequest {
	return CreateLancamentoRequest{
		IDLoja:      1,
		IDComanda:   2,
		IDAtendente: 3,
		DataHora:    time.Now().Format(time.RFC3339),
		Finalizado:  boolPtr(false),
	}
}

func TestServiceCreate(t *testing.T) {
	t.Run("validation failures", func(t *testing.T) {
		cases := []CreateLancamentoRequest{
			{IDComanda: 1, IDAtendente: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: boolPtr(false)},
			{IDLoja: 1, IDAtendente: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: boolPtr(false)},
			{IDLoja: 1, IDComanda: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: boolPtr(false)},
			{IDLoja: 1, IDComanda: 1, IDAtendente: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: nil},
			{IDLoja: 1, IDComanda: 1, IDAtendente: 1, DataHora: "", Finalizado: boolPtr(false)},
			{IDLoja: 1, IDComanda: 1, IDAtendente: 1, DataHora: "invalid-date", Finalizado: boolPtr(false)},
		}

		svc := NewService(repositoryStub{})
		for i, req := range cases {
			_, err := svc.Create(context.Background(), req)
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("case %d expected ErrValidation, got %v", i, err)
			}
		}
	})

	t.Run("returns duplicate and repository errors", func(t *testing.T) {
		repoErr := errors.New("repo")
		req := validCreateRequest()

		svcErr := NewService(repositoryStub{existsByLojaComandaFn: func(_ context.Context, _, _ int, _ bool) (bool, error) {
			return false, repoErr
		}})
		_, err := svcErr.Create(context.Background(), req)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected repo error, got %v", err)
		}

		svcDup := NewService(repositoryStub{existsByLojaComandaFn: func(_ context.Context, _, _ int, _ bool) (bool, error) {
			return true, nil
		}})
		_, err = svcDup.Create(context.Background(), req)
		if !errors.Is(err, ErrDuplicateLancamento) {
			t.Fatalf("expected ErrDuplicateLancamento, got %v", err)
		}
	})

	t.Run("create persistence error and success", func(t *testing.T) {
		req := validCreateRequest()
		repoErr := errors.New("create failed")

		svcErr := NewService(repositoryStub{createFn: func(_ context.Context, _ *models.LancamentoComanda) error {
			return repoErr
		}})
		_, err := svcErr.Create(context.Background(), req)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected create error, got %v", err)
		}

		svcOK := NewService(repositoryStub{})
		model, err := svcOK.Create(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if model == nil || model.IDLoja != req.IDLoja || model.IDComanda != req.IDComanda || model.IDAtendente != req.IDAtendente || model.Finalizado != *req.Finalizado {
			t.Fatalf("unexpected model: %+v", model)
		}
	})
}

func TestServiceList(t *testing.T) {
	t.Run("invalid filters", func(t *testing.T) {
		svc := NewService(repositoryStub{})
		_, err := svc.List(context.Background(), ListLancamentosRequest{IDComanda: "x"})
		if !errors.Is(err, ErrInvalidFilter) {
			t.Fatalf("expected ErrInvalidFilter for IDComanda, got %v", err)
		}
		_, err = svc.List(context.Background(), ListLancamentosRequest{IDMesa: "x"})
		if !errors.Is(err, ErrInvalidFilter) {
			t.Fatalf("expected ErrInvalidFilter for IDMesa, got %v", err)
		}
		_, err = svc.List(context.Background(), ListLancamentosRequest{IDAtendente: "x"})
		if !errors.Is(err, ErrInvalidFilter) {
			t.Fatalf("expected ErrInvalidFilter for IDAtendente, got %v", err)
		}
		_, err = svc.List(context.Background(), ListLancamentosRequest{DataHora: "x"})
		if !errors.Is(err, ErrInvalidFilter) {
			t.Fatalf("expected ErrInvalidFilter for DataHora, got %v", err)
		}
	})

	t.Run("repository error and success with parsed filters", func(t *testing.T) {
		repoErr := errors.New("list failed")
		called := 0
		svcErr := NewService(repositoryStub{listFn: func(_ context.Context, _ ListLancamentosFilter) ([]models.LancamentoComanda, error) {
			called++
			return nil, repoErr
		}})
		_, err := svcErr.List(context.Background(), ListLancamentosRequest{IDComanda: "10"})
		if !errors.Is(err, repoErr) || called != 1 {
			t.Fatalf("expected repo error once, got err=%v called=%d", err, called)
		}

		var captured ListLancamentosFilter
		svcOK := NewService(repositoryStub{listFn: func(_ context.Context, f ListLancamentosFilter) ([]models.LancamentoComanda, error) {
			captured = f
			return []models.LancamentoComanda{{ID: 1}}, nil
		}})
		finalizado := true
		result, err := svcOK.List(context.Background(), ListLancamentosRequest{IDComanda: "1", IDMesa: "2", IDAtendente: "3", DataHora: "2026-01-02", Finalizado: &finalizado})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || captured.IDComanda == nil || *captured.IDComanda != 1 || captured.IDMesa == nil || *captured.IDMesa != 2 || captured.IDAtendente == nil || *captured.IDAtendente != 3 || captured.DataHora == nil || captured.Finalizado == nil || *captured.Finalizado != true {
			t.Fatalf("unexpected captured filter/result: %+v %+v", captured, result)
		}
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("validation failures", func(t *testing.T) {
		cases := []CreateLancamentoRequest{
			{IDComanda: 1, IDAtendente: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: boolPtr(false)},
			{IDLoja: 1, IDAtendente: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: boolPtr(false)},
			{IDLoja: 1, IDComanda: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: boolPtr(false)},
			{IDLoja: 1, IDComanda: 1, IDAtendente: 1, DataHora: time.Now().Format(time.RFC3339), Finalizado: nil},
			{IDLoja: 1, IDComanda: 1, IDAtendente: 1, DataHora: "", Finalizado: boolPtr(false)},
			{IDLoja: 1, IDComanda: 1, IDAtendente: 1, DataHora: "invalid-date", Finalizado: boolPtr(false)},
		}

		svc := NewService(repositoryStub{})
		for i, req := range cases {
			_, err := svc.Update(context.Background(), 1, req)
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("case %d expected ErrValidation, got %v", i, err)
			}
		}
	})

	t.Run("find not found and generic errors", func(t *testing.T) {
		req := validCreateRequest()
		svcNF := NewService(repositoryStub{findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
			return nil, gorm.ErrRecordNotFound
		}})
		_, err := svcNF.Update(context.Background(), 1, req)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}

		repoErr := errors.New("find failed")
		svcErr := NewService(repositoryStub{findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
			return nil, repoErr
		}})
		_, err = svcErr.Update(context.Background(), 1, req)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected repo error, got %v", err)
		}
	})

	t.Run("finalizado true to false validation", func(t *testing.T) {
		req := validCreateRequest()
		svc := NewService(repositoryStub{findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
			return &models.LancamentoComanda{ID: 1, Finalizado: true}, nil
		}})
		_, err := svc.Update(context.Background(), 1, req)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation, got %v", err)
		}
	})

	t.Run("duplicate checks and persistence", func(t *testing.T) {
		req := validCreateRequest()
		initial := &models.LancamentoComanda{ID: 1, Finalizado: false}

		repoErr := errors.New("exists failed")
		svcExistsErr := NewService(repositoryStub{
			findByIDFn:                       func(_ context.Context, _ uint) (*models.LancamentoComanda, error) { return initial, nil },
			existsByLojaComandaExcludingIDFn: func(_ context.Context, _ uint, _, _ int, _ bool) (bool, error) { return false, repoErr },
		})
		_, err := svcExistsErr.Update(context.Background(), 1, req)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected exists error, got %v", err)
		}

		svcDup := NewService(repositoryStub{
			findByIDFn:                       func(_ context.Context, _ uint) (*models.LancamentoComanda, error) { return initial, nil },
			existsByLojaComandaExcludingIDFn: func(_ context.Context, _ uint, _, _ int, _ bool) (bool, error) { return true, nil },
		})
		_, err = svcDup.Update(context.Background(), 1, req)
		if !errors.Is(err, ErrDuplicateLancamento) {
			t.Fatalf("expected ErrDuplicateLancamento, got %v", err)
		}

		updateErr := errors.New("update failed")
		svcUpdateErr := NewService(repositoryStub{
			findByIDFn:                       func(_ context.Context, _ uint) (*models.LancamentoComanda, error) { return initial, nil },
			existsByLojaComandaExcludingIDFn: func(_ context.Context, _ uint, _, _ int, _ bool) (bool, error) { return false, nil },
			updateFn:                         func(_ context.Context, _ *models.LancamentoComanda) error { return updateErr },
		})
		_, err = svcUpdateErr.Update(context.Background(), 1, req)
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected update error, got %v", err)
		}

		svcOK := NewService(repositoryStub{
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 1}, nil
			},
			existsByLojaComandaExcludingIDFn: func(_ context.Context, _ uint, _, _ int, _ bool) (bool, error) { return false, nil },
		})
		model, err := svcOK.Update(context.Background(), 1, req)
		if err != nil || model == nil || model.IDComanda != req.IDComanda {
			t.Fatalf("unexpected update result: model=%+v err=%v", model, err)
		}
	})
}

func TestServiceCreateItems(t *testing.T) {
	makeReq := func() CreateItemsRequest {
		return CreateItemsRequest{Itens: []CreateItemRequest{{IDLancamentoComanda: 1, Sequencia: 1, IDProduto: 2, Quantidade: 1, PrecoVenda: 10, IDAtendente: 3}}}
	}

	t.Run("validation failures", func(t *testing.T) {
		req := makeReq()
		req.Itens[0].Quantidade = 0
		svc := NewService(repositoryStub{})
		_, err := svc.CreateItems(context.Background(), req)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation for quantidade, got %v", err)
		}

		req = makeReq()
		req.Itens[0].PrecoVenda = -1
		_, err = svc.CreateItems(context.Background(), req)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation for precoVenda, got %v", err)
		}
	})

	t.Run("find lancamento errors and finalizado", func(t *testing.T) {
		req := makeReq()
		svcNF := NewService(repositoryStub{findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
			return nil, gorm.ErrRecordNotFound
		}})
		_, err := svcNF.CreateItems(context.Background(), req)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}

		repoErr := errors.New("find failed")
		svcErr := NewService(repositoryStub{findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
			return nil, repoErr
		}})
		_, err = svcErr.CreateItems(context.Background(), req)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected repo error, got %v", err)
		}

		svcFin := NewService(repositoryStub{findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
			return &models.LancamentoComanda{ID: 1, Finalizado: true}, nil
		}})
		_, err = svcFin.CreateItems(context.Background(), req)
		if !errors.Is(err, ErrComandaFinalizada) {
			t.Fatalf("expected ErrComandaFinalizada, got %v", err)
		}
	})

	t.Run("sequencia checks and create batch", func(t *testing.T) {
		req := makeReq()
		repoErr := errors.New("sequencia check failed")
		svcSeqErr := NewService(repositoryStub{
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 1}, nil
			},
			sequenciaExistsInLancamentoFn: func(_ context.Context, _ uint, _ int) (bool, error) { return false, repoErr },
		})
		_, err := svcSeqErr.CreateItems(context.Background(), req)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected sequencia error, got %v", err)
		}

		svcDup := NewService(repositoryStub{
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 1}, nil
			},
			sequenciaExistsInLancamentoFn: func(_ context.Context, _ uint, _ int) (bool, error) { return true, nil },
		})
		_, err = svcDup.CreateItems(context.Background(), req)
		if !errors.Is(err, ErrDuplicateSequencia) {
			t.Fatalf("expected ErrDuplicateSequencia, got %v", err)
		}

		reqBatchDup := CreateItemsRequest{Itens: []CreateItemRequest{{IDLancamentoComanda: 1, Sequencia: 1, IDProduto: 1, Quantidade: 1, IDAtendente: 1}, {IDLancamentoComanda: 1, Sequencia: 1, IDProduto: 2, Quantidade: 1, IDAtendente: 1}}}
		svcBatchDup := NewService(repositoryStub{
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 1}, nil
			},
			sequenciaExistsInLancamentoFn: func(_ context.Context, _ uint, _ int) (bool, error) { return false, nil },
		})
		_, err = svcBatchDup.CreateItems(context.Background(), reqBatchDup)
		if !errors.Is(err, ErrDuplicateSequencia) {
			t.Fatalf("expected ErrDuplicateSequencia (batch dup), got %v", err)
		}

		batchErr := errors.New("batch failed")
		svcCreateErr := NewService(repositoryStub{
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 1}, nil
			},
			sequenciaExistsInLancamentoFn: func(_ context.Context, _ uint, _ int) (bool, error) { return false, nil },
			createItemsBatchFn:            func(_ context.Context, _ []*models.LancamentoComandaItem) error { return batchErr },
		})
		_, err = svcCreateErr.CreateItems(context.Background(), req)
		if !errors.Is(err, batchErr) {
			t.Fatalf("expected batch error, got %v", err)
		}

		svcOK := NewService(repositoryStub{
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 1}, nil
			},
			sequenciaExistsInLancamentoFn: func(_ context.Context, _ uint, _ int) (bool, error) { return false, nil },
		})
		items, err := svcOK.CreateItems(context.Background(), req)
		if err != nil || len(items) != 1 {
			t.Fatalf("unexpected success result: items=%+v err=%v", items, err)
		}
	})
}

func TestServiceUpdateItem(t *testing.T) {
	t.Run("validation failures", func(t *testing.T) {
		svc := NewService(repositoryStub{})
		_, err := svc.UpdateItem(context.Background(), 1, UpdateItemRequest{})
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation for empty request, got %v", err)
		}
		_, err = svc.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(0)})
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation for quantidade <= 0, got %v", err)
		}
	})

	t.Run("find errors and finalizado", func(t *testing.T) {
		svcItemNF := NewService(repositoryStub{findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) {
			return nil, gorm.ErrRecordNotFound
		}})
		_, err := svcItemNF.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(1)})
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound for item, got %v", err)
		}

		repoErr := errors.New("find item failed")
		svcItemErr := NewService(repositoryStub{findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) { return nil, repoErr }})
		_, err = svcItemErr.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(1)})
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected find item error, got %v", err)
		}

		baseItem := &models.LancamentoComandaItem{ID: 1, IDLancamentoComanda: 10}
		svcLancNF := NewService(repositoryStub{
			findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) { return baseItem, nil },
			findByIDFn:     func(_ context.Context, _ uint) (*models.LancamentoComanda, error) { return nil, gorm.ErrRecordNotFound },
		})
		_, err = svcLancNF.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(1)})
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound for lancamento, got %v", err)
		}

		repoFindErr := errors.New("find lanc failed")
		svcLancErr := NewService(repositoryStub{
			findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) { return baseItem, nil },
			findByIDFn:     func(_ context.Context, _ uint) (*models.LancamentoComanda, error) { return nil, repoFindErr },
		})
		_, err = svcLancErr.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(1)})
		if !errors.Is(err, repoFindErr) {
			t.Fatalf("expected find lanc error, got %v", err)
		}

		svcFin := NewService(repositoryStub{
			findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) { return baseItem, nil },
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 10, Finalizado: true}, nil
			},
		})
		_, err = svcFin.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(1)})
		if !errors.Is(err, ErrComandaFinalizada) {
			t.Fatalf("expected ErrComandaFinalizada, got %v", err)
		}
	})

	t.Run("update persistence and success mutations", func(t *testing.T) {
		baseItem := &models.LancamentoComandaItem{ID: 1, IDLancamentoComanda: 10, Quantidade: 1}
		updateErr := errors.New("update failed")
		svcUpdateErr := NewService(repositoryStub{
			findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) {
				copy := *baseItem
				return &copy, nil
			},
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 10, Finalizado: false}, nil
			},
			updateItemFn: func(_ context.Context, _ *models.LancamentoComandaItem) error { return updateErr },
		})
		_, err := svcUpdateErr.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(5)})
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected update error, got %v", err)
		}

		svcQty := NewService(repositoryStub{
			findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) {
				copy := *baseItem
				return &copy, nil
			},
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 10, Finalizado: false}, nil
			},
		})
		item, err := svcQty.UpdateItem(context.Background(), 1, UpdateItemRequest{Quantidade: floatPtr(5)})
		if err != nil || item.Quantidade != 5 {
			t.Fatalf("expected quantidade update, got item=%+v err=%v", item, err)
		}

		svcCancelTrue := NewService(repositoryStub{
			findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) {
				return &models.LancamentoComandaItem{ID: 1, IDLancamentoComanda: 10}, nil
			},
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 10, Finalizado: false}, nil
			},
		})
		item, err = svcCancelTrue.UpdateItem(context.Background(), 1, UpdateItemRequest{Cancelado: boolPtr(true)})
		if err != nil || !item.Cancelado || item.DataHoraCancelamento == nil {
			t.Fatalf("expected cancel true with timestamp, got item=%+v err=%v", item, err)
		}

		cancelTime := time.Now()
		svcCancelFalse := NewService(repositoryStub{
			findItemByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComandaItem, error) {
				return &models.LancamentoComandaItem{ID: 1, IDLancamentoComanda: 10, Cancelado: true, DataHoraCancelamento: &cancelTime}, nil
			},
			findByIDFn: func(_ context.Context, _ uint) (*models.LancamentoComanda, error) {
				return &models.LancamentoComanda{ID: 10, Finalizado: false}, nil
			},
		})
		item, err = svcCancelFalse.UpdateItem(context.Background(), 1, UpdateItemRequest{Cancelado: boolPtr(false)})
		if err != nil || item.Cancelado || item.DataHoraCancelamento != nil {
			t.Fatalf("expected cancel false without timestamp, got item=%+v err=%v", item, err)
		}
	})
}

func TestServiceListItens(t *testing.T) {
	t.Run("invalid id_comanda", func(t *testing.T) {
		svc := NewService(repositoryStub{})
		_, err := svc.ListItens(context.Background(), ListItensRequest{IDComanda: "abc"})
		if !errors.Is(err, ErrInvalidFilter) {
			t.Fatalf("expected ErrInvalidFilter, got %v", err)
		}
	})

	t.Run("repository error and mapping success", func(t *testing.T) {
		repoErr := errors.New("query failed")
		svcErr := NewService(repositoryStub{listItensByComandaFn: func(_ context.Context, _ int) ([]ItemComandaRow, error) {
			return nil, repoErr
		}})
		_, err := svcErr.ListItens(context.Background(), ListItensRequest{IDComanda: "1"})
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected repo error, got %v", err)
		}

		now := time.Now()
		svcOK := NewService(repositoryStub{listItensByComandaFn: func(_ context.Context, id int) ([]ItemComandaRow, error) {
			if id != 1 {
				t.Fatalf("expected id 1, got %d", id)
			}
			return []ItemComandaRow{{LancamentoComandaItem: models.LancamentoComandaItem{ID: 7, IDLancamentoComanda: 8, Sequencia: 1, IDProduto: 2, CodigoBarras: "123", Quantidade: 3, PrecoVenda: 9.5, Cancelado: true, DataHoraCancelamento: &now, IDAtendente: 4, IDSituacao: 5}, IDComanda: 99}}, nil
		}})

		result, err := svcOK.ListItens(context.Background(), ListItensRequest{IDComanda: "1"})
		if err != nil || len(result) != 1 || result[0].ID != 7 || result[0].IDComanda != 99 {
			t.Fatalf("unexpected list itens result: %+v err=%v", result, err)
		}
	})
}

func TestParseDataHora(t *testing.T) {
	if _, err := parseDataHora(time.Now().Format(time.RFC3339)); err != nil {
		t.Fatalf("expected parse success, got %v", err)
	}
	if _, err := parseDataHora("2026-01-02 15:04:05"); err != nil {
		t.Fatalf("expected parse success, got %v", err)
	}
	if _, err := parseDataHora("2026-01-02"); err != nil {
		t.Fatalf("expected parse success, got %v", err)
	}
	if _, err := parseDataHora("invalid"); err == nil {
		t.Fatal("expected parse error for invalid date")
	}
}
