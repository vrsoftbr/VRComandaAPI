package parametros

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router gin.IRouter, getDatabase func() *mongo.Database, invalidateConnection func()) {
	repository := NewMongoRepository(getDatabase, invalidateConnection, "parametros")
	handler := NewHandler(NewServiceWithRepository(repository))

	router.GET("/parametros", handler.List)
}
