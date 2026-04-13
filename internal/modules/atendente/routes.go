package atendente

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes wires Atendente HTTP endpoints and its dependencies.
func RegisterRoutes(router *gin.Engine, getDatabase func() *mongo.Database, invalidateConnection func()) {
	repository := NewMongoRepository(getDatabase, invalidateConnection, "atendentes")
	service := NewService(repository)
	handler := NewHandler(service)

	router.GET("/atendentes", handler.List)
}
