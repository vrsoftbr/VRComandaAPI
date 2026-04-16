package mesa

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var originalFindAllMesasFn = findAllMesasFn

func resetFindAllMesasFn() {
	findAllMesasFn = originalFindAllMesasFn
}

func TestMongoRepositoryListBuildsQueryAndSort(t *testing.T) {
	defer resetFindAllMesasFn()

	expected := []Mesa{{Mesa: 1, Descricao: "Descricao", Ativo: true}}

	findAllMesasFn = func(ctx context.Context, collection *mongo.Collection, filter interface{}, findOptions *options.FindOptions, invalidateConnection func()) ([]Mesa, error) {
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
		collectionName: "mesas",
	}

	ativo := true
	ctx := context.WithValue(context.Background(), "k", "v")
	result, err := repo.List(ctx, ListMesasFilter{
		IDLoja: 10,
		Mesa:   1,
		Ativo:  &ativo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Mesa != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestMongoRepositoryListWithEmptyFilter(t *testing.T) {
	defer resetFindAllMesasFn()

	var capturedFilter interface{}
	findAllMesasFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]Mesa, error) {
		capturedFilter = filter
		return []Mesa{}, nil
	}

	repo := &mongoRepository{
		getDatabase:          func() *mongo.Database { return nil },
		invalidateConnection: func() {},
		collectionName:       "mesas",
	}

	result, err := repo.List(context.Background(), ListMesasFilter{})
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
