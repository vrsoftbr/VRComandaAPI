package comanda

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
// @Summary Listar comandas
// @Tags Comanda
// @Produce json
// @Param id_loja query int false "ID da loja"
// @Param comanda query int false "Numero da comanda"
// @Param numero_identificacao query string false "Numero de identificacao"
// @Param ativo query bool false "Ativo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /comandas [get]
func (h *Handler) List(c *gin.Context) {
	var req ListComandasRequest
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
