package lancamento

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/shared/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Update godoc
// @Summary Editar lancamento
// @Tags Lancamento
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body CreateLancamentoRequest true "Lancamento"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lancamentos/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "id invalido")
		return
	}

	var req CreateLancamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrDuplicateLancamento) {
			utils.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}

// CreateItem godoc
// @Summary Criar itens de lancamento em lote
// @Tags Lancamento
// @Accept json
// @Produce json
// @Param body body CreateItemsRequest true "Itens"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lancamentos/itens [post]
func (h *Handler) CreateItem(c *gin.Context) {
	var req CreateItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.CreateItems(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrDuplicateSequencia) {
			utils.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, ErrComandaFinalizada) {
			utils.RespondError(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusCreated, result)
}

// UpdateItem godoc
// @Summary Editar item de lancamento
// @Tags Lancamento
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body UpdateItemRequest true "Item"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lancamentos/itens/{id} [put]
func (h *Handler) UpdateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "id invalido")
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.UpdateItem(c.Request.Context(), uint(id), req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrComandaFinalizada) {
			utils.RespondError(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}

// ListItens godoc
// @Summary Listar itens por comanda
// @Tags Lancamento
// @Produce json
// @Param id_comanda query int true "ID comanda"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lancamentos/itens [get]
func (h *Handler) ListItens(c *gin.Context) {
	var req ListItensRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.ListItens(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidFilter) {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}

// Create godoc
// @Summary Criar lancamento
// @Tags Lancamento
// @Accept json
// @Produce json
// @Param body body CreateLancamentoRequest true "Lancamento"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lancamentos [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateLancamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrDuplicateLancamento) {
			utils.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusCreated, result)
}

// List godoc
// @Summary Listar lancamentos
// @Tags Lancamento
// @Produce json
// @Param id_comanda query int false "ID comanda"
// @Param id_mesa query int false "ID mesa"
// @Param id_atendente query int false "ID atendente"
// @Param dataHora query string false "DataHora"
// @Param finalizado query bool false "Finalizado"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lancamentos [get]
func (h *Handler) List(c *gin.Context) {
	var req ListLancamentosRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidFilter) {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}
