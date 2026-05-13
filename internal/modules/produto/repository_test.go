package produto

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var originalFindAllProdutosCodigoBarrasFn = findAllProdutosCodigoBarrasFn
var originalFindAllProdutosFn = findAllProdutosFn
var originalCountProdutosFn = countProdutosFn

func resetFindFns() {
	findAllProdutosCodigoBarrasFn = originalFindAllProdutosCodigoBarrasFn
	findAllProdutosFn = originalFindAllProdutosFn
	countProdutosFn = originalCountProdutosFn
}

func TestMongoRepositoryListUsesCodigoBarrasToFilterProdutos(t *testing.T) {
	defer resetFindFns()

	callsCodigoBarras := 0
	var codigoBarrasFilter interface{}
	findAllProdutosCodigoBarrasFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]ProdutoCodigoBarras, error) {
		callsCodigoBarras++
		if callsCodigoBarras == 1 {
			codigoBarrasFilter = filter
			return []ProdutoCodigoBarras{{IDProduto: 1}, {IDProduto: 1}, {IDProduto: 2}}, nil
		}
		return []ProdutoCodigoBarras{{IDProduto: 1, CodigoBarras: "12345"}, {IDProduto: 2, CodigoBarras: "12345"}}, nil
	}

	var produtosFilter interface{}
	findAllProdutosFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]Produto, error) {
		produtosFilter = filter
		return []Produto{{IDProduto: 1}, {IDProduto: 2}}, nil
	}

	repo := &mongoRepository{
		getDatabase:                func() *mongo.Database { return nil },
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	result, err := repo.List(context.Background(), ListProdutosFilter{
		IDLoja:            10,
		CodigoBarras:      " 12345 ",
		DescricaoCompleta: "REPOLHO",
		DescricaoCupom:    "KG",
		Page:              1,
		Limit:             20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("expected 2 produtos, got %d", len(result.Items))
	}

	cbQuery, ok := codigoBarrasFilter.(bson.M)
	if !ok {
		t.Fatalf("expected bson.M for codigo barras query, got %T", codigoBarrasFilter)
	}
	orQuery, ok := cbQuery["$or"].([]bson.M)
	if !ok || len(orQuery) != 2 {
		t.Fatalf("expected $or with two barcode fields, got %+v", cbQuery["$or"])
	}

	prodQuery, ok := produtosFilter.(bson.M)
	if !ok {
		t.Fatalf("expected bson.M for produtos query, got %T", produtosFilter)
	}
	if prodQuery["idLoja"] != 10 {
		t.Fatalf("unexpected idLoja in query: %+v", prodQuery["idLoja"])
	}

	descCompleta, ok := prodQuery["descricaocompleta"].(bson.M)
	if !ok {
		t.Fatalf("expected descricaocompleta regex filter, got %T", prodQuery["descricaocompleta"])
	}
	if descCompleta["$options"] != "i" {
		t.Fatalf("expected case-insensitive regex option, got %+v", descCompleta["$options"])
	}

	descCupom, ok := prodQuery["descricaocupom"].(bson.M)
	if !ok {
		t.Fatalf("expected descricaocupom regex filter, got %T", prodQuery["descricaocupom"])
	}
	if descCupom["$options"] != "i" {
		t.Fatalf("expected case-insensitive regex option, got %+v", descCupom["$options"])
	}

	idProdutoFilter, ok := prodQuery["idProduto"].(bson.M)
	if !ok {
		t.Fatalf("expected idProduto bson.M filter, got %T", prodQuery["idProduto"])
	}
	ids, ok := idProdutoFilter["$in"].([]int)
	if !ok {
		t.Fatalf("expected $in []int, got %T", idProdutoFilter["$in"])
	}
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 2 {
		t.Fatalf("unexpected idProduto $in values: %+v", ids)
	}
	if callsCodigoBarras != 2 {
		t.Fatalf("expected 2 barcode queries, got %d", callsCodigoBarras)
	}
	if len(result.Items[0].CodigosBarras) == 0 {
		t.Fatal("expected barcode data attached to products")
	}
}

func TestMongoRepositoryListReturnsEmptyWhenCodigoBarrasNotFound(t *testing.T) {
	defer resetFindFns()

	findAllProdutosCodigoBarrasFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]ProdutoCodigoBarras, error) {
		return []ProdutoCodigoBarras{}, nil
	}

	calledProdutosQuery := false
	findAllProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]Produto, error) {
		calledProdutosQuery = true
		return []Produto{{IDProduto: 1}}, nil
	}

	repo := &mongoRepository{
		getDatabase:                func() *mongo.Database { return nil },
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	result, err := repo.List(context.Background(), ListProdutosFilter{CodigoBarras: "999", Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Items) != 0 {
		t.Fatalf("expected empty list, got len=%d", len(result.Items))
	}
	if calledProdutosQuery {
		t.Fatal("produtos query should not run when barcode lookup returns empty")
	}
}

func TestMongoRepositoryListWithoutCodigoBarrasQueriesProdutosDirectly(t *testing.T) {
	defer resetFindFns()

	calledCodigoBarrasQuery := false
	findAllProdutosCodigoBarrasFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]ProdutoCodigoBarras, error) {
		calledCodigoBarrasQuery = true
		return nil, nil
	}

	var capturedFilter interface{}
	findAllProdutosFn = func(_ context.Context, _ *mongo.Collection, filter interface{}, _ *options.FindOptions, _ func()) ([]Produto, error) {
		capturedFilter = filter
		return []Produto{}, nil
	}

	repo := &mongoRepository{
		getDatabase:                func() *mongo.Database { return nil },
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	_, err := repo.List(context.Background(), ListProdutosFilter{IDLoja: 1, Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledCodigoBarrasQuery {
		t.Fatal("codigoBarras query should not run when filter has no codigoBarras")
	}

	query, ok := capturedFilter.(bson.M)
	if !ok {
		t.Fatalf("expected bson.M filter, got %T", capturedFilter)
	}
	if query["idLoja"] != 1 {
		t.Fatalf("unexpected idLoja in query: %+v", query["idLoja"])
	}
}

func TestMongoRepositoryCollectionReturnsNilWhenDatabaseNil(t *testing.T) {
	repo := &mongoRepository{
		getDatabase:                func() *mongo.Database { return nil },
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	if repo.collection("produtos") != nil {
		t.Fatal("expected nil collection when DB is nil")
	}
}

func TestMongoRepositoryCountDocumentsErrorPropagated(t *testing.T) {
	defer resetFindFns()

	countProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}) (int64, error) {
		return 0, errors.New("count failed")
	}

	fakeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://fake:27017"))
	if err != nil {
		t.Fatalf("unexpected client creation error: %v", err)
	}

	repo := &mongoRepository{
		getDatabase: func() *mongo.Database {
			return fakeClient.Database("vrcomanda_test")
		},
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	_, err = repo.List(context.Background(), ListProdutosFilter{IDLoja: 1, Page: 1, Limit: 20})
	if err == nil {
		t.Fatal("expected error from countProdutosFn, got nil")
	}
}

func TestMongoRepositoryCountZeroReturnsEmptyResult(t *testing.T) {
	defer resetFindFns()

	countProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}) (int64, error) {
		return 0, nil // count returns 0 → should return empty result immediately
	}

	fakeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://fake:27017"))
	if err != nil {
		t.Fatalf("unexpected client creation error: %v", err)
	}

	repo := &mongoRepository{
		getDatabase: func() *mongo.Database {
			return fakeClient.Database("vrcomanda_test")
		},
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	result, err := repo.List(context.Background(), ListProdutosFilter{IDLoja: 1, Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Items) != 0 || result.Total != 0 {
		t.Fatalf("expected empty result when count == 0, got %+v", result)
	}
}

func TestMongoRepositoryCountPositiveProceedsToPaginate(t *testing.T) {
	defer resetFindFns()

	countProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}) (int64, error) {
		return 3, nil // count returns 3 → proceed to findAll
	}

	findAllProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]Produto, error) {
		return []Produto{{IDProduto: 1, IDLoja: 1}}, nil
	}

	findAllProdutosCodigoBarrasFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]ProdutoCodigoBarras, error) {
		return []ProdutoCodigoBarras{}, nil
	}

	fakeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://fake:27017"))
	if err != nil {
		t.Fatalf("unexpected client creation error: %v", err)
	}

	repo := &mongoRepository{
		getDatabase: func() *mongo.Database {
			return fakeClient.Database("vrcomanda_test")
		},
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	result, err := repo.List(context.Background(), ListProdutosFilter{IDLoja: 1, Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Items) != 1 || result.Total != 3 {
		t.Fatalf("expected one item with total=3, got %+v", result)
	}
}

func TestMongoRepositoryFindAllProdutosErrorPropagated(t *testing.T) {
	defer resetFindFns()

	countProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}) (int64, error) {
		return 5, nil
	}

	findAllProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]Produto, error) {
		return nil, errors.New("findAll products failed")
	}

	fakeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://fake:27017"))
	if err != nil {
		t.Fatalf("unexpected client creation error: %v", err)
	}

	repo := &mongoRepository{
		getDatabase: func() *mongo.Database {
			return fakeClient.Database("vrcomanda_test")
		},
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	_, err = repo.List(context.Background(), ListProdutosFilter{IDLoja: 1, Page: 1, Limit: 20})
	if err == nil {
		t.Fatal("expected error from findAllProdutosFn, got nil")
	}
}

func TestMongoRepositoryCodigoBarrasLookupErrorPropagated(t *testing.T) {
	defer resetFindFns()

	// When CodigoBarras filter is empty, there's only one call: the product ID barcode lookup
	findAllProdutosCodigoBarrasFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]ProdutoCodigoBarras, error) {
		return nil, errors.New("codigo barras lookup failed")
	}

	findAllProdutosFn = func(_ context.Context, _ *mongo.Collection, _ interface{}, _ *options.FindOptions, _ func()) ([]Produto, error) {
		return []Produto{{IDProduto: 1, IDLoja: 1}}, nil
	}

	repo := &mongoRepository{
		getDatabase:                func() *mongo.Database { return nil },
		invalidateConnection:       func() {},
		produtosCollectionName:     "produtos",
		codigoBarrasCollectionName: "produtoscodigobarras",
	}

	_, err := repo.List(context.Background(), ListProdutosFilter{IDLoja: 1, Page: 1, Limit: 20})
	if err == nil {
		t.Fatal("expected error from findAllProdutosCodigoBarrasFn product lookup, got nil")
	}
}
