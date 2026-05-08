package produto

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
// @Summary Listar produtos
// @Tags Produto
// @Produce json
// @Param idLoja query int false "ID da loja"
// @Param codigoBarras query string false "Codigo de barras"
// @Param descricaocompleta query string false "Descricao completa"
// @Param descricaocupom query string false "Descricao cupom"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /produtos [get]
func (h *Handler) List(c *gin.Context) {
	var req ListProdutosRequest
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
