package mesa

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes wires Mesa HTTP endpoints and its dependencies.
func RegisterRoutes(router *gin.Engine, getDatabase func() *mongo.Database, invalidateConnection func()) {
	repository := NewMongoRepository(getDatabase, invalidateConnection, "mesas")
	service := NewService(repository)
	handler := NewHandler(service)

	router.GET("/mesas", handler.List)
}
