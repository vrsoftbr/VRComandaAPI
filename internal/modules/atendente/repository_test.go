package atendente

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var originalFindAllAtendentesFn = findAllAtendentesFn

func resetFindAllAtendentesFn() {
	findAllAtendentesFn = originalFindAllAtendentesFn
}

func TestMongoRepositoryListBuildsQueryAndSort(t *testing.T) {
	defer resetFindAllAtendentesFn()

	expected := []Atendente{{Codigo: "A1", Nome: "Ana", Ativo: true}}

	findAllAtendentesFn = func(ctx context.Context, collection *mongo.Collection, filter interface{}, findOptions *options.FindOptions, invalidateConnection func()) ([]Atendente, error) {
		return expected, nil
	}

	invalidateCalls := 0
	fakeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://fake:27017"))
	if err != nil {
		t.Fatalf("unexpected client creation error: %v", err)
	}

	repo := &mongoRepository{
		getDatabase: func() *mongo.Database {
			return fakeClient.Database("vrcomanda_test")
		},
		invalidateConnection: func() {
			invalidateCalls++
		},
		collectionName: "atendentes",
	}

	ativo := true
	ctx := context.WithValue(context.Background(), "k", "v")
	result, err := repo.List(ctx, ListAtendentesFilter{
		IDLoja: 10,
		Codigo: "A1",
		Nome:   "An(a)",
		Ativo:  &ativo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Codigo != "A1" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestMongoRepositoryListWithEmptyFilter(t *testing.T) {
	defer resetFindAllAtendentesFn()

	var capturedFilter interface{}
	findAllAtendentesFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]Atendente, error) {
		capturedFilter = filter
		return []Atendente{}, nil
	}

	repo := &mongoRepository{
		getDatabase:          func() *mongo.Database { return nil },
		invalidateConnection: func() {},
		collectionName:       "atendentes",
	}

	result, err := repo.List(context.Background(), ListAtendentesFilter{})
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
