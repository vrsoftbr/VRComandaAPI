package lancamento

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vrcomandaapi/internal/shared/models"
)

func newRepositoryForTest(t *testing.T, migrate bool) Repository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if migrate {
		if err := db.AutoMigrate(&models.LancamentoComanda{}, &models.LancamentoComandaItem{}); err != nil {
			t.Fatalf("failed to migrate sqlite: %v", err)
		}
	}

	return NewRepository(db)
}

func seedLancamento(t *testing.T, repo Repository, idLoja, idComanda int, finalizado bool) *models.LancamentoComanda {
	t.Helper()

	finalizadoVal := finalizado
	model := &models.LancamentoComanda{
		IDLoja:      idLoja,
		IDComanda:   idComanda,
		IDAtendente: 1,
		DataHora:    time.Now(),
		Finalizado:  finalizadoVal,
	}
	if err := repo.Create(context.Background(), model); err != nil {
		t.Fatalf("seed create failed: %v", err)
	}
	return model
}

func TestRepositoryCreateAndFindByID(t *testing.T) {
	repo := newRepositoryForTest(t, true)
	model := seedLancamento(t, repo, 1, 10, false)

	found, err := repo.FindByID(context.Background(), model.ID)
	if err != nil {
		t.Fatalf("unexpected find error: %v", err)
	}
	if found.ID != model.ID || found.IDComanda != 10 {
		t.Fatalf("unexpected found model: %+v", found)
	}
}

func TestRepositoryCreateReturnsErrorWithoutTables(t *testing.T) {
	repo := newRepositoryForTest(t, false)
	err := repo.Create(context.Background(), &models.LancamentoComanda{IDLoja: 1, IDComanda: 1, IDAtendente: 1, DataHora: time.Now(), Finalizado: false})
	if err == nil {
		t.Fatal("expected create error without migrated tables")
	}
}

func TestRepositoryExistsMethods(t *testing.T) {
	repo := newRepositoryForTest(t, true)
	a := seedLancamento(t, repo, 1, 11, false)
	seedLancamento(t, repo, 1, 11, true)

	exists, err := repo.ExistsByLojaComanda(context.Background(), 1, 11, false)
	if err != nil || !exists {
		t.Fatalf("expected exists=true for matching record, exists=%v err=%v", exists, err)
	}
	exists, err = repo.ExistsByLojaComanda(context.Background(), 1, 999, false)
	if err != nil || exists {
		t.Fatalf("expected exists=false for missing record, exists=%v err=%v", exists, err)
	}

	exists, err = repo.ExistsByLojaComandaExcludingID(context.Background(), a.ID, 1, 11, false)
	if err != nil || exists {
		t.Fatalf("expected exists=false excluding only matching row, exists=%v err=%v", exists, err)
	}

	seedLancamento(t, repo, 1, 11, false)
	exists, err = repo.ExistsByLojaComandaExcludingID(context.Background(), a.ID, 1, 11, false)
	if err != nil || !exists {
		t.Fatalf("expected exists=true with another row, exists=%v err=%v", exists, err)
	}
}

func TestRepositoryExistsMethodsReturnErrorWithoutTables(t *testing.T) {
	repo := newRepositoryForTest(t, false)
	if _, err := repo.ExistsByLojaComanda(context.Background(), 1, 1, false); err == nil {
		t.Fatal("expected ExistsByLojaComanda error without tables")
	}
	if _, err := repo.ExistsByLojaComandaExcludingID(context.Background(), 1, 1, 1, false); err == nil {
		t.Fatal("expected ExistsByLojaComandaExcludingID error without tables")
	}
}

func TestRepositoryUpdateAndList(t *testing.T) {
	repo := newRepositoryForTest(t, true)
	model := seedLancamento(t, repo, 2, 20, false)
	idMesa := 5
	model.IDMesa = &idMesa
	model.Observacao = "updated"
	model.Finalizado = true
	model.DataHora = time.Now().Add(1 * time.Hour)
	if err := repo.Update(context.Background(), model); err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}

	idComanda := 20
	idAtendente := 1
	finalizado := true
	result, err := repo.List(context.Background(), ListLancamentosFilter{
		IDComanda:   &idComanda,
		IDMesa:      &idMesa,
		IDAtendente: &idAtendente,
		Finalizado:  &finalizado,
		DataHora:    &model.DataHora,
	})
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}
	if len(result) != 1 || result[0].Observacao != "updated" {
		t.Fatalf("unexpected list result: %+v", result)
	}
}

func TestRepositoryUpdateAndListReturnErrorsWithoutTables(t *testing.T) {
	repo := newRepositoryForTest(t, false)
	if err := repo.Update(context.Background(), &models.LancamentoComanda{ID: 1}); err == nil {
		t.Fatal("expected update error without tables")
	}
	if _, err := repo.List(context.Background(), ListLancamentosFilter{}); err == nil {
		t.Fatal("expected list error without tables")
	}
}

func TestRepositoryFindByIDAndFindItemByIDNotFound(t *testing.T) {
	repo := newRepositoryForTest(t, true)
	if _, err := repo.FindByID(context.Background(), 999); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound for FindByID, got %v", err)
	}
	if _, err := repo.FindItemByID(context.Background(), 999); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound for FindItemByID, got %v", err)
	}
}

func TestRepositoryItemsMethods(t *testing.T) {
	repo := newRepositoryForTest(t, true)
	lanc := seedLancamento(t, repo, 3, 30, false)

	items := []*models.LancamentoComandaItem{{IDLancamentoComanda: lanc.ID, Sequencia: 1, IDProduto: 10, Quantidade: 2, PrecoVenda: 5, IDAtendente: 1, IDSituacao: 1}}
	if err := repo.CreateItemsBatch(context.Background(), items); err != nil {
		t.Fatalf("unexpected CreateItemsBatch error: %v", err)
	}
	if items[0].ID == 0 {
		t.Fatal("expected created item with non-zero ID")
	}

	exists, err := repo.SequenciaExistsInLancamento(context.Background(), lanc.ID, 1)
	if err != nil || !exists {
		t.Fatalf("expected sequencia exists=true, exists=%v err=%v", exists, err)
	}
	exists, err = repo.SequenciaExistsInLancamento(context.Background(), lanc.ID, 2)
	if err != nil || exists {
		t.Fatalf("expected sequencia exists=false, exists=%v err=%v", exists, err)
	}

	found, err := repo.FindItemByID(context.Background(), items[0].ID)
	if err != nil || found.ID != items[0].ID {
		t.Fatalf("unexpected find item result: item=%+v err=%v", found, err)
	}

	found.Quantidade = 9
	if err := repo.UpdateItem(context.Background(), found); err != nil {
		t.Fatalf("unexpected update item error: %v", err)
	}

	rows, err := repo.ListItensByComanda(context.Background(), 30)
	if err != nil {
		t.Fatalf("unexpected ListItensByComanda error: %v", err)
	}
	if len(rows) != 1 || rows[0].IDComanda != 30 || rows[0].Quantidade != 9 {
		t.Fatalf("unexpected rows: %+v", rows)
	}
}

func TestRepositoryItemsMethodsReturnErrorsWithoutTables(t *testing.T) {
	repo := newRepositoryForTest(t, false)
	if err := repo.CreateItemsBatch(context.Background(), []*models.LancamentoComandaItem{{IDLancamentoComanda: 1, Sequencia: 1, IDProduto: 1, Quantidade: 1, IDAtendente: 1}}); err == nil {
		t.Fatal("expected CreateItemsBatch error without tables")
	}
	if _, err := repo.SequenciaExistsInLancamento(context.Background(), 1, 1); err == nil {
		t.Fatal("expected SequenciaExistsInLancamento error without tables")
	}
	if err := repo.UpdateItem(context.Background(), &models.LancamentoComandaItem{ID: 1}); err == nil {
		t.Fatal("expected UpdateItem error without tables")
	}
	if _, err := repo.ListItensByComanda(context.Background(), 1); err == nil {
		t.Fatal("expected ListItensByComanda error without tables")
	}
}
