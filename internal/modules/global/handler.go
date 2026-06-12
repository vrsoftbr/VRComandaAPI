package global

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/shared/utils"
)

type Handler struct {
	lancamentosDetalhesService LancamentosDetalhesService
	consultarSituacaoService   ComandaCatracaService
	comandaPDVService          ComandaPDVService
}

func NewHandler(lancamentosDetalhesService LancamentosDetalhesService, consultarSituacaoService ComandaCatracaService, comandaPDVService ComandaPDVService) *Handler {
	return &Handler{
		lancamentosDetalhesService: lancamentosDetalhesService,
		consultarSituacaoService:   consultarSituacaoService,
		comandaPDVService:          comandaPDVService,
	}
}

// GetLancamentosDetalhes godoc
// @Summary Listar lancamentos detalhados
// @Description Lista lancamentos com detalhes de comanda, mesa e itens.
// @Tags Global
// @Param id_loja query int false "ID da loja"
// @Param finalizado query bool false "Filtrar por status finalizado"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lancamentos/detalhes [get]
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

// ComandaCatraca godoc
// @Summary Consultar situacao de comanda
// @Description Consulta situacao da comanda por numero de identificacao.
// @Tags Global
// @Param idLoja query int true "ID da loja"
// @Param numeroIdentificacaoComanda query string true "Numero de identificacao da comanda"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /comanda/consultarsituacao [get]
func (h *Handler) ComandaCatraca(c *gin.Context) {
	var req ComandaCatracaRequest
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

// GetLancamentosPDV godoc
// @Summary Consultar comanda para o PDV
// @Description Consulta a comanda aberta e seus itens para venda no PDV.
// @Tags Global
// @Param numeroComanda query int true "Numero da comanda"
// @Param loja query int true "ID da loja"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /venda/comanda/pdv/consultar [get]
func (h *Handler) GetLancamentosPDV(c *gin.Context) {
	var req GetLancamentoPDVRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondMappedError(c, fmt.Errorf("%w: %v", ErrInvalidRequest, err), errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest})
		return
	}

	result, err := h.comandaPDVService.Consultar(c.Request.Context(), req)
	if err != nil {
		respondMappedError(c, err,
			errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest},
			errorMapping{target: lancamento.ErrInvalidFilter, status: http.StatusBadRequest},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "mensagem": "Comanda"})
}

// AtualizarComandaPDV godoc
// @Summary Atualizar situacao da comanda
// @Description Atualiza somente o campo finalizado do lancamento da comanda.
// @Tags Global
// @Accept json
// @Produce json
// @Param body body UpdadeLancamentoPDVRequest true "Situacao da comanda"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /atualizacomanda [put]
func (h *Handler) UpdateComandaPDV(c *gin.Context) {
	var req UpdadeLancamentoPDVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondMappedError(c, fmt.Errorf("%w: %v", ErrInvalidRequest, err), errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest})
		return
	}

	err := h.comandaPDVService.Atualizar(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, lancamento.ErrNotFound) {
			utils.RespondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondMappedError(c, err,
			errorMapping{target: ErrInvalidRequest, status: http.StatusBadRequest},
			errorMapping{target: lancamento.ErrValidation, status: http.StatusBadRequest},
		)
		return
	}

	utils.RespondOK(c, http.StatusOK, gin.H{"mensagem": "Comanda atualizada com sucesso"})
}
