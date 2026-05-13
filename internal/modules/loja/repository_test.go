package loja

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var originalFindAllLojasFn = findAllLojasFn

func resetFindAllLojasFn() {
	findAllLojasFn = originalFindAllLojasFn
}

func TestMongoRepositoryListBuildsEmptyQuery(t *testing.T) {
	defer resetFindAllLojasFn()

	var capturedFilter interface{}
	findAllLojasFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]Loja, error) {
		capturedFilter = filter
		return []Loja{{ID: 1, Descricao: "Loja A"}}, nil
	}

	repo := &mongoRepository{
		getDatabase:          func() *mongo.Database { return nil },
		invalidateConnection: func() {},
		collectionName:       "lojas",
	}

	result, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].ID != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}

	query, ok := capturedFilter.(bson.M)
	if !ok {
		t.Fatalf("expected bson.M filter, got %T", capturedFilter)
	}
	if len(query) != 0 {
		t.Fatalf("expected empty filter, got %+v", query)
	}
}

func TestMongoRepositoryListReturnsEmptyWhenCollectionNil(t *testing.T) {
	defer resetFindAllLojasFn()

	findAllLojasFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]Loja, error) {
		return []Loja{}, nil
	}

	repo := &mongoRepository{
		getDatabase:          func() *mongo.Database { return nil },
		invalidateConnection: func() {},
		collectionName:       "lojas",
	}

	result, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got len=%d", len(result))
	}
}

func TestMongoRepositoryCollectionReturnsNilWhenDatabaseNil(t *testing.T) {
	repo := &mongoRepository{
		getDatabase:          func() *mongo.Database { return nil },
		invalidateConnection: func() {},
		collectionName:       "lojas",
	}

	if repo.collection() != nil {
		t.Fatal("expected nil collection when DB is nil")
	}
}

func TestMongoRepositoryCollectionReturnsCollectionWhenDatabaseNotNil(t *testing.T) {
	fakeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://fake:27017"))
	if err != nil {
		t.Fatalf("unexpected client creation error: %v", err)
	}

	repo := &mongoRepository{
		getDatabase: func() *mongo.Database {
			return fakeClient.Database("vrcomanda_test")
		},
		invalidateConnection: func() {},
		collectionName:       "lojas",
	}

	if repo.collection() == nil {
		t.Fatal("expected non-nil collection when DB is not nil")
	}
}
