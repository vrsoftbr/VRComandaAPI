package comanda

import (
	"context"

	"vrcomandaapi/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	List(ctx context.Context, filter ListComandasFilter) ([]Comanda, error)
}

var findAllComandasFn = database.FindAll[Comanda]

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

func (r *mongoRepository) List(ctx context.Context, filter ListComandasFilter) ([]Comanda, error) {
	collection := r.collection()

	query := bson.M{}
	if filter.IDLoja != 0 {
		query["idLoja"] = filter.IDLoja
	}
	if filter.Comanda != 0 {
		query["comanda"] = filter.Comanda
	}
	if filter.NumeroIdentificacao != "" {
		query["numeroIdentificacao"] = filter.NumeroIdentificacao
	}
	if filter.Ativo != nil {
		query["ativo"] = *filter.Ativo
	}

	findOptions := options.Find().SetSort(bson.M{"comanda": 1})
	return findAllComandasFn(ctx, collection, query, findOptions, r.invalidateConnection)
}
