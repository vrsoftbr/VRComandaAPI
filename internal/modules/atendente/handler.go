package atendente

import (
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
// @Summary Listar atendentes
// @Tags Atendente
// @Produce json
// @Param id_loja query int false "ID da loja"
// @Param codigo query string false "Codigo"
// @Param nome query string false "Nome"
// @Param ativo query bool false "Ativo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /atendentes [get]
func (h *Handler) List(c *gin.Context) {
	var req ListAtendentesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}
