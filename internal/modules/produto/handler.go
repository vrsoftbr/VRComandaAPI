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
// @Param page query int false "Número da página (padrão: 1)"
// @Param limit query int false "Itens por página (padrão: 20, máximo: 100)"
// @Success 200 {object} map[string]interface{}  "items, page, limit, total, pages"
// @Failure 400 {object} map[string]interface{}  "Parâmetros inválidos"
// @Failure 500 {object} map[string]interface{}  "Erro interno do servidor"
// @Router /produtos [get]
func (h *Handler) List(c *gin.Context) {
	var req ListProdutosRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Configurar paginação com valores padrão
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}
