package parametros

import (
	"context"

	"vrcomandaapi/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	ListByLojaAndParametros(ctx context.Context, idLoja int, ids []int) ([]Parametro, error)
}

var findAllParametrosFn = database.FindAll[Parametro]

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

func (r *mongoRepository) ListByLojaAndParametros(ctx context.Context, idLoja int, ids []int) ([]Parametro, error) {
	query := bson.M{
		"idLoja": idLoja,
		"idParametro": bson.M{
			"$in": ids,
		},
	}

	findOptions := options.Find().SetSort(bson.M{"idParametro": 1})
	return findAllParametrosFn(ctx, r.collection(), query, findOptions, r.invalidateConnection)
}
