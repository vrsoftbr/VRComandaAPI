package global

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/utils"
)

type Handler struct {
	lancamentosDetalhesService LancamentosDetalhesService
}

func NewHandler(lancamentosDetalhesService LancamentosDetalhesService) *Handler {
	return &Handler{
		lancamentosDetalhesService: lancamentosDetalhesService,
	}
}

func (h *Handler) GetLancamentosDetalhes(c *gin.Context) {
	var req ListLancamentosDetalhesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondMappedError(c, fmt.Errorf("%w: %v", ErrInvalidRequest, err), errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest})
		return
	}

	result, err := h.lancamentosDetalhesService.Execute(c.Request.Context(), req)
	if err != nil {
		respondMappedError(c, err,
			errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest},
			errorMapping{target: lancamento.ErrInvalidFilter, status: http.StatusBadRequest},
		)
		return
	}

	utils.RespondOK(c, http.StatusOK, result)
}
