package loja

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes wires Loja HTTP endpoints and its dependencies.
func RegisterRoutes(router gin.IRouter, getDatabase func() *mongo.Database, invalidateConnection func()) {
	repository := NewMongoRepository(getDatabase, invalidateConnection, "lojas")
	service := NewService(repository)
	handler := NewHandler(service)

	router.GET("/lojas", handler.List)
}
