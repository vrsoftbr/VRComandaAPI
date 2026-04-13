package mesa

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
// @Summary Listar mesas
// @Tags Mesa
// @Produce json
// @Param id_loja query int false "ID da loja"
// @Param mesa query int false "Numero da mesa"
// @Param ativo query bool false "Ativo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mesas [get]
func (h *Handler) List(c *gin.Context) {
	var req ListMesasRequest
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
