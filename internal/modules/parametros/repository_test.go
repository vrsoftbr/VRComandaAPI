package parametros

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var originalFindAllParametrosFn = findAllParametrosFn

func resetFindAllParametrosFn() {
	findAllParametrosFn = originalFindAllParametrosFn
}

func TestMongoRepositoryListByLojaAndParametrosBuildsQuery(t *testing.T) {
	defer resetFindAllParametrosFn()

	var capturedFilter interface{}
	findAllParametrosFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]Parametro, error) {
		capturedFilter = filter
		return []Parametro{{IDParametro: 14, IDLoja: 1, Valor: "2"}}, nil
	}

	repo := &mongoRepository{
		getDatabase:          func() *mongo.Database { return nil },
		invalidateConnection: func() {},
		collectionName:       "parametros",
	}

	result, err := repo.ListByLojaAndParametros(context.Background(), 1, []int{13, 14, 15})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].IDParametro != 14 {
		t.Fatalf("unexpected result: %+v", result)
	}

	query, ok := capturedFilter.(bson.M)
	if !ok {
		t.Fatalf("expected bson.M filter, got %T", capturedFilter)
	}
	if query["idLoja"] != 1 {
		t.Fatalf("unexpected idLoja in query: %+v", query["idLoja"])
	}
	idParametroFilter, ok := query["idParametro"].(bson.M)
	if !ok {
		t.Fatalf("expected bson.M idParametro filter, got %T", query["idParametro"])
	}
	ids, ok := idParametroFilter["$in"].([]int)
	if !ok {
		t.Fatalf("expected $in []int, got %T", idParametroFilter["$in"])
	}
	if len(ids) != 3 || ids[0] != 13 || ids[1] != 14 || ids[2] != 15 {
		t.Fatalf("unexpected idParametro $in values: %+v", ids)
	}
}
