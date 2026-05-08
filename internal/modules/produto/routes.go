package produto

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes wires Produto HTTP endpoints and its dependencies.
func RegisterRoutes(router gin.IRouter, getDatabase func() *mongo.Database, invalidateConnection func()) {
	repository := NewMongoRepository(getDatabase, invalidateConnection, "produtos", "produtoscodigobarras")
	service := NewService(repository)
	handler := NewHandler(service)

	router.GET("/produtos", handler.List)
}
