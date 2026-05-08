package global

import (
	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
)

func RegisterRoutes(router gin.IRouter, lancamentoService lancamento.Service, comandaService comanda.Service, mesaService mesa.Service) {
	service := NewLancamentosDetalhesService(lancamentoService, comandaService, mesaService)
	consultarSituacaoService := NewConsultarSituacaoComandaService(lancamentoService, comandaService)
	h := NewHandler(service, consultarSituacaoService)
	router.GET("/lancamentos/detalhes", h.GetLancamentosDetalhes)
	router.GET("/comanda/consultarsituacao", h.ConsultarSituacaoComanda)
}
