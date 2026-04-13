package lancamento

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes wires Lancamento HTTP endpoints and its dependencies.
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	router.POST("/lancamentos", handler.Create)
	router.PUT("/lancamentos/:id", handler.Update)
	router.GET("/lancamentos", handler.List)
	router.POST("/lancamentos/itens", handler.CreateItem)
	router.PUT("/lancamentos/itens/:id", handler.UpdateItem)
	router.GET("/lancamentos/itens", handler.ListItens)
}
