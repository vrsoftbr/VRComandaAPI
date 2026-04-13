package comanda

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes wires Comanda HTTP endpoints and its dependencies.
func RegisterRoutes(router *gin.Engine, getDatabase func() *mongo.Database, invalidateConnection func()) {
	repository := NewMongoRepository(getDatabase, invalidateConnection, "comandas")
	service := NewService(repository)
	handler := NewHandler(service)

	router.GET("/comandas", handler.List)
}
