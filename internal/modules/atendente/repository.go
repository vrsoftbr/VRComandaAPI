package atendente

import (
	"context"
	"regexp"

	"vrcomandaapi/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	List(ctx context.Context, filter ListAtendentesFilter) ([]Atendente, error)
}

var findAllAtendentesFn = database.FindAll[Atendente]

type mongoRepository struct {
	getDatabase          func() *mongo.Database
	invalidateConnection func()
	collectionName       string
}

func NewMongoRepository(getDatabase func() *mongo.Database, invalidateConnection func(), collectionName string) Repository {
	return &mongoRepository{getDatabase: getDatabase, invalidateConnection: invalidateConnection, collectionName: collectionName}
}

func (r *mongoRepository) collection() *mongo.Collection {
	db := r.getDatabase()
	if db == nil {
		return nil
	}

	return db.Collection(r.collectionName)
}

func (r *mongoRepository) List(ctx context.Context, filter ListAtendentesFilter) ([]Atendente, error) {
	collection := r.collection()

	query := bson.M{}
	if filter.IDLoja != 0 {
		query["idLoja"] = filter.IDLoja
	}
	if filter.Codigo != "" {
		query["codigo"] = filter.Codigo
	}
	if filter.Nome != "" {
		pattern := ".*" + regexp.QuoteMeta(filter.Nome) + ".*"
		query["nome"] = bson.M{"$regex": pattern, "$options": "i"}
	}
	if filter.Ativo != nil {
		query["ativo"] = *filter.Ativo
	}

	findOptions := options.Find().SetSort(bson.M{"nome": 1})
	return findAllAtendentesFn(ctx, collection, query, findOptions, r.invalidateConnection)
}
