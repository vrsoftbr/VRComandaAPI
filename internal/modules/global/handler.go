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
	consultarSituacaoService   ConsultarSituacaoComandaService
}

func NewHandler(lancamentosDetalhesService LancamentosDetalhesService, consultarSituacaoService ConsultarSituacaoComandaService) *Handler {
	return &Handler{
		lancamentosDetalhesService: lancamentosDetalhesService,
		consultarSituacaoService:   consultarSituacaoService,
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

func (h *Handler) ConsultarSituacaoComanda(c *gin.Context) {
	var req ConsultarSituacaoComandaRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondMappedError(c, fmt.Errorf("%w: %v", ErrInvalidRequest, err), errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest})
		return
	}

	result, err := h.consultarSituacaoService.Execute(c.Request.Context(), req)
	if err != nil {
		respondMappedError(c, err,
			errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest},
		)
		return
	}

	if result == nil {
		c.JSON(http.StatusOK, gin.H{"data": nil, "mensagem": "Comanda não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "mensagem": "Comanda encontrada"})
}
