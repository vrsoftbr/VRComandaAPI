package parametros

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/shared/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// List godoc
// @Summary Listar parametros da loja
// @Tags Parametros
// @Produce json
// @Param idLoja query int true "ID da loja"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /parametros [get]
func (h *Handler) List(c *gin.Context) {
	var req ListParametrosRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, fmt.Errorf("%w: %v", ErrInvalidRequest, err).Error())
		return
	}

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrInvalidRequest) {
			status = http.StatusBadRequest
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}
