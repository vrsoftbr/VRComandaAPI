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
	if collection == nil {
		return []Comanda{}, nil
	}

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

	cursor, err := collection.Find(ctx, query, options.Find().SetSort(bson.M{"comanda": 1}))
	if err != nil {
		if database.IsMongoConnectionError(err) {
			r.invalidateConnection()
		}
		return []Comanda{}, nil
	}
	defer cursor.Close(ctx)

	var result []Comanda
	if err := cursor.All(ctx, &result); err != nil {
		if database.IsMongoConnectionError(err) {
			r.invalidateConnection()
		}
		return []Comanda{}, nil
	}

	if result == nil {
		return []Comanda{}, nil
	}

	return result, nil
}
