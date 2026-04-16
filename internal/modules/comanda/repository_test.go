package comanda

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var originalFindAllComandasFn = findAllComandasFn

func resetFindAllComandasFn() {
	findAllComandasFn = originalFindAllComandasFn
}

func TestMongoRepositoryListBuildsQueryAndSort(t *testing.T) {
	defer resetFindAllComandasFn()

	expected := []Comanda{{Comanda: 1, NumeroIdentificacao: "X1", Ativo: true}}

	findAllComandasFn = func(ctx context.Context, collection *mongo.Collection, filter interface{}, findOptions *options.FindOptions, invalidateConnection func()) ([]Comanda, error) {
		return expected, nil
	}

	invalideCalls := 0
	fakeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://fake:27017"))
	if err != nil {
		t.Fatalf("unexpected client creation error: %v", err)
	}

	repo := &mongoRepository{
		getDatabase: func() *mongo.Database {
			return fakeClient.Database("vrcomanda_test")
		},
		invalidateConnection: func() {
			invalideCalls++
		},
		collectionName: "comandas",
	}

	ativo := true
	ctx := context.WithValue(context.Background(), "k", "v")
	result, err := repo.List(ctx, ListComandasFilter{
		IDLoja:              10,
		Comanda:             1,
		NumeroIdentificacao: "X1",
		Ativo:               &ativo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Comanda != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestMongoRepositoryListWithEmptyFilter(t *testing.T) {
	defer resetFindAllComandasFn()

	var capturedFilter interface{}
	findAllComandasFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]Comanda, error) {
		capturedFilter = filter
		return []Comanda{}, nil
	}

	repo := &mongoRepository{
		getDatabase:          func() *mongo.Database { return nil },
		invalidateConnection: func() {},
		collectionName:       "comandas",
	}

	result, err := repo.List(context.Background(), ListComandasFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := capturedFilter.(bson.M)
	if !ok {
		t.Fatalf("expected bson.M filter, got %T", capturedFilter)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got len=%d", len(result))
	}
}
