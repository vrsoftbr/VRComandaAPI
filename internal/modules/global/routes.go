package global

import (
	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
)

func RegisterRoutes(router *gin.Engine, lancamentoService lancamento.Service, comandaService comanda.Service, mesaService mesa.Service) {
	service := NewLancamentosDetalhesService(lancamentoService, comandaService, mesaService)
	h := NewHandler(service)
	router.GET("/lancamentos/detalhes", h.GetLancamentosDetalhes)
}
