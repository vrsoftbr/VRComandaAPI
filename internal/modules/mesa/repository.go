package mesa

import (
	"context"

	"vrcomandaapi/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	List(ctx context.Context, filter ListMesasFilter) ([]Mesa, error)
}

var findAllMesasFn = database.FindAll[Mesa]

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

func (r *mongoRepository) List(ctx context.Context, filter ListMesasFilter) ([]Mesa, error) {
	collection := r.collection()

	query := bson.M{}
	if filter.IDLoja != 0 {
		query["idLoja"] = filter.IDLoja
	}
	if filter.Mesa != 0 {
		query["mesa"] = filter.Mesa
	}
	if len(filter.Mesas) > 0 {
		query["mesa"] = bson.M{"$in": filter.Mesas}
	}
	if filter.Ativo != nil {
		query["ativo"] = *filter.Ativo
	}

	findOptions := options.Find().SetSort(bson.M{"mesa": 1})
	return findAllMesasFn(ctx, collection, query, findOptions, r.invalidateConnection)
}
