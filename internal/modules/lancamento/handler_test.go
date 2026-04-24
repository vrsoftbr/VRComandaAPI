package lancamento

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"vrcomandaapi/internal/shared/models"
	"vrcomandaapi/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

// serviceStub implements Service with configurable function fields.
// Fields not set default to returning zero values without error.
type serviceStub struct {
	createFn      func(ctx context.Context, req CreateLancamentoRequest) (*models.LancamentoComanda, error)
	updateFn      func(ctx context.Context, id uint, req CreateLancamentoRequest) (*models.LancamentoComanda, error)
	listFn        func(ctx context.Context, req ListLancamentosRequest) ([]models.LancamentoComanda, error)
	listItensFn   func(ctx context.Context, req ListItensRequest) ([]ItemComandaResponse, error)
	createItemsFn func(ctx context.Context, req CreateItemsRequest) ([]*models.LancamentoComandaItem, error)
	updateItemFn  func(ctx context.Context, id uint, req UpdateItemRequest) (*models.LancamentoComandaItem, error)
}

func (s serviceStub) Create(ctx context.Context, req CreateLancamentoRequest) (*models.LancamentoComanda, error) {
	if s.createFn != nil {
		return s.createFn(ctx, req)
	}
	return nil, nil
}

func (s serviceStub) Update(ctx context.Context, id uint, req CreateLancamentoRequest) (*models.LancamentoComanda, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, req)
	}
	return nil, nil
}

func (s serviceStub) List(ctx context.Context, req ListLancamentosRequest) ([]models.LancamentoComanda, error) {
	if s.listFn != nil {
		return s.listFn(ctx, req)
	}
	return []models.LancamentoComanda{}, nil
}

func (s serviceStub) ListItens(ctx context.Context, req ListItensRequest) ([]ItemComandaResponse, error) {
	if s.listItensFn != nil {
		return s.listItensFn(ctx, req)
	}
	return []ItemComandaResponse{}, nil
}

func (s serviceStub) CreateItems(ctx context.Context, req CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
	if s.createItemsFn != nil {
		return s.createItemsFn(ctx, req)
	}
	return nil, nil
}

func (s serviceStub) UpdateItem(ctx context.Context, id uint, req UpdateItemRequest) (*models.LancamentoComandaItem, error) {
	if s.updateItemFn != nil {
		return s.updateItemFn(ctx, id, req)
	}
	return nil, nil
}

func newRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(svc)
	r.POST("/lancamentos", h.Create)
	r.PUT("/lancamentos/:id", h.Update)
	r.GET("/lancamentos", h.List)
	r.POST("/lancamentos/itens", h.CreateItem)
	r.PUT("/lancamentos/itens/:id", h.UpdateItem)
	r.GET("/lancamentos/itens", h.ListItens)
	return r
}

// Create
func TestHandlerCreate(t *testing.T) {
	t.Run("returns 400 when body is malformed JSON", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertDataNil(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
	})

	t.Run("returns 400 when service returns ErrValidation", func(t *testing.T) {
		r := newRouter(serviceStub{createFn: func(_ context.Context, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return nil, ErrValidation
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertDataNil(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
	})

	t.Run("returns 409 when service returns ErrDuplicateLancamento", func(t *testing.T) {
		r := newRouter(serviceStub{createFn: func(_ context.Context, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return nil, ErrDuplicateLancamento
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 on generic service error", func(t *testing.T) {
		r := newRouter(serviceStub{createFn: func(_ context.Context, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return nil, errors.New("db error")
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 201 on success", func(t *testing.T) {
		r := newRouter(serviceStub{createFn: func(_ context.Context, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return &models.LancamentoComanda{ID: 1}, nil
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertMessageEquals(t, utils.DecodeBodyMap(t, w.Body.Bytes()), "ok")
	})
}

// Update
func TestHandlerUpdate(t *testing.T) {
	t.Run("returns 400 when id param is not numeric", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/abc", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when body is malformed JSON", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/1", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when service returns ErrValidation", func(t *testing.T) {
		r := newRouter(serviceStub{updateFn: func(_ context.Context, _ uint, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return nil, ErrValidation
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 409 when service returns ErrDuplicateLancamento", func(t *testing.T) {
		r := newRouter(serviceStub{updateFn: func(_ context.Context, _ uint, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return nil, ErrDuplicateLancamento
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 404 when service returns ErrNotFound", func(t *testing.T) {
		r := newRouter(serviceStub{updateFn: func(_ context.Context, _ uint, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return nil, ErrNotFound
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 on generic service error", func(t *testing.T) {
		r := newRouter(serviceStub{updateFn: func(_ context.Context, _ uint, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return nil, errors.New("db error")
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 on success", func(t *testing.T) {
		r := newRouter(serviceStub{updateFn: func(_ context.Context, _ uint, _ CreateLancamentoRequest) (*models.LancamentoComanda, error) {
			return &models.LancamentoComanda{ID: 1}, nil
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertMessageEquals(t, utils.DecodeBodyMap(t, w.Body.Bytes()), "ok")
	})
}

// List
func TestHandlerList(t *testing.T) {
	t.Run("returns 400 when query bind fails", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos?finalizado=not-bool", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertDataNil(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
	})

	t.Run("returns 400 when service returns ErrInvalidFilter", func(t *testing.T) {
		r := newRouter(serviceStub{listFn: func(_ context.Context, _ ListLancamentosRequest) ([]models.LancamentoComanda, error) {
			return nil, ErrInvalidFilter
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 on generic service error", func(t *testing.T) {
		r := newRouter(serviceStub{listFn: func(_ context.Context, _ ListLancamentosRequest) ([]models.LancamentoComanda, error) {
			return nil, errors.New("db error")
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 on success", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		data := utils.AssertDataArray(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
		if len(data) != 0 {
			t.Fatalf("expected empty data, got len=%d", len(data))
		}
	})
}

// CreateItem
func TestHandlerCreateItem(t *testing.T) {
	t.Run("returns 400 when body is malformed JSON", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertDataNil(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
	})

	t.Run("returns 400 when service returns ErrValidation", func(t *testing.T) {
		r := newRouter(serviceStub{createItemsFn: func(_ context.Context, _ CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
			return nil, ErrValidation
		}})

		w := httptest.NewRecorder()
		body := `{"itens":[{"id_lancamentocomanda":1,"sequencia":1,"id_produto":1,"quantidade":1,"id_atendente":1}]}`
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 404 when service returns ErrNotFound", func(t *testing.T) {
		r := newRouter(serviceStub{createItemsFn: func(_ context.Context, _ CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
			return nil, ErrNotFound
		}})

		w := httptest.NewRecorder()
		body := `{"itens":[{"id_lancamentocomanda":1,"sequencia":1,"id_produto":1,"quantidade":1,"id_atendente":1}]}`
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 409 when service returns ErrDuplicateSequencia", func(t *testing.T) {
		r := newRouter(serviceStub{createItemsFn: func(_ context.Context, _ CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
			return nil, ErrDuplicateSequencia
		}})

		w := httptest.NewRecorder()
		body := `{"itens":[{"id_lancamentocomanda":1,"sequencia":1,"id_produto":1,"quantidade":1,"id_atendente":1}]}`
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 422 when service returns ErrComandaFinalizada", func(t *testing.T) {
		r := newRouter(serviceStub{createItemsFn: func(_ context.Context, _ CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
			return nil, ErrComandaFinalizada
		}})

		w := httptest.NewRecorder()
		body := `{"itens":[{"id_lancamentocomanda":1,"sequencia":1,"id_produto":1,"quantidade":1,"id_atendente":1}]}`
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 on generic service error", func(t *testing.T) {
		r := newRouter(serviceStub{createItemsFn: func(_ context.Context, _ CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
			return nil, errors.New("db error")
		}})

		w := httptest.NewRecorder()
		body := `{"itens":[{"id_lancamentocomanda":1,"sequencia":1,"id_produto":1,"quantidade":1,"id_atendente":1}]}`
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 201 on success", func(t *testing.T) {
		r := newRouter(serviceStub{createItemsFn: func(_ context.Context, _ CreateItemsRequest) ([]*models.LancamentoComandaItem, error) {
			return []*models.LancamentoComandaItem{{ID: 1}}, nil
		}})

		w := httptest.NewRecorder()
		body := `{"itens":[{"id_lancamentocomanda":1,"sequencia":1,"id_produto":1,"quantidade":1,"id_atendente":1}]}`
		req := httptest.NewRequest(http.MethodPost, "/lancamentos/itens", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertMessageEquals(t, utils.DecodeBodyMap(t, w.Body.Bytes()), "ok")
	})
}

// UpdateItem
func TestHandlerUpdateItem(t *testing.T) {
	t.Run("returns 400 when id param is not numeric", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/abc", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when body is malformed JSON", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/1", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 400 when service returns ErrValidation", func(t *testing.T) {
		r := newRouter(serviceStub{updateItemFn: func(_ context.Context, _ uint, _ UpdateItemRequest) (*models.LancamentoComandaItem, error) {
			return nil, ErrValidation
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 404 when service returns ErrNotFound", func(t *testing.T) {
		r := newRouter(serviceStub{updateItemFn: func(_ context.Context, _ uint, _ UpdateItemRequest) (*models.LancamentoComandaItem, error) {
			return nil, ErrNotFound
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 422 when service returns ErrComandaFinalizada", func(t *testing.T) {
		r := newRouter(serviceStub{updateItemFn: func(_ context.Context, _ uint, _ UpdateItemRequest) (*models.LancamentoComandaItem, error) {
			return nil, ErrComandaFinalizada
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 on generic service error", func(t *testing.T) {
		r := newRouter(serviceStub{updateItemFn: func(_ context.Context, _ uint, _ UpdateItemRequest) (*models.LancamentoComandaItem, error) {
			return nil, errors.New("db error")
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 on success", func(t *testing.T) {
		r := newRouter(serviceStub{updateItemFn: func(_ context.Context, _ uint, _ UpdateItemRequest) (*models.LancamentoComandaItem, error) {
			return &models.LancamentoComandaItem{ID: 1}, nil
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lancamentos/itens/1", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertMessageEquals(t, utils.DecodeBodyMap(t, w.Body.Bytes()), "ok")
	})
}

// ListItens
func TestHandlerListItens(t *testing.T) {
	t.Run("returns 400 when id_comanda is missing", func(t *testing.T) {
		r := newRouter(serviceStub{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/itens", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
		utils.AssertDataNil(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
	})

	t.Run("returns 400 when service returns ErrInvalidFilter", func(t *testing.T) {
		r := newRouter(serviceStub{listItensFn: func(_ context.Context, _ ListItensRequest) ([]ItemComandaResponse, error) {
			return nil, ErrInvalidFilter
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/itens?id_comanda=5", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 500 on generic service error", func(t *testing.T) {
		r := newRouter(serviceStub{listItensFn: func(_ context.Context, _ ListItensRequest) ([]ItemComandaResponse, error) {
			return nil, errors.New("db error")
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/itens?id_comanda=5", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("returns 200 on success", func(t *testing.T) {
		r := newRouter(serviceStub{listItensFn: func(_ context.Context, req ListItensRequest) ([]ItemComandaResponse, error) {
			if req.IDComanda != 5 {
				t.Fatalf("unexpected id_comanda: %d", req.IDComanda)
			}
			return []ItemComandaResponse{{ID: 1}}, nil
		}})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lancamentos/itens?id_comanda=5", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		data := utils.AssertDataArray(t, utils.DecodeBodyMap(t, w.Body.Bytes()))
		if len(data) != 1 {
			t.Fatalf("expected one item, got %d", len(data))
		}
	})
}
