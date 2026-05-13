package global

import (
	"github.com/gin-gonic/gin"

	"vrcomandaapi/internal/modules/atendente"
	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/lancamento"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/modules/produto"
)

func RegisterRoutes(router gin.IRouter, lancamentoService lancamento.Service, atendenteService atendente.Service, comandaService comanda.Service, mesaService mesa.Service, produtoService produto.Service) {
	service := NewLancamentosDetalhesService(lancamentoService, atendenteService, comandaService, mesaService, produtoService)
	consultarSituacaoService := NewConsultarSituacaoComandaService(lancamentoService, comandaService)
	h := NewHandler(service, consultarSituacaoService)
	router.GET("/lancamentos/detalhes", h.GetLancamentosDetalhes)
	router.GET("/comanda/consultarsituacao", h.ConsultarSituacaoComanda)
}
