package loja

import (
	"context"

	"vrcomandaapi/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	List(ctx context.Context) ([]Loja, error)
}

var findAllLojasFn = database.FindAll[Loja]

type mongoRepository struct {
	getDatabase          func() *mongo.Database
	invalidateConnection func()
	collectionName       string
}

func NewMongoRepository(
	getDatabase func() *mongo.Database,
	invalidateConnection func(),
	collectionName string,
) Repository {
	return &mongoRepository{
		getDatabase:          getDatabase,
		invalidateConnection: invalidateConnection,
		collectionName:       collectionName,
	}
}

func (r *mongoRepository) collection() *mongo.Collection {
	db := r.getDatabase()
	if db == nil {
		return nil
	}
	return db.Collection(r.collectionName)
}

func (r *mongoRepository) List(ctx context.Context) ([]Loja, error) {
	items, err := findAllLojasFn(ctx, r.collection(), bson.M{}, nil, r.invalidateConnection)
	if err != nil {
		return nil, err
	}
	return items, nil
}
