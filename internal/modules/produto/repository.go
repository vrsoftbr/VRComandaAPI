package produto

import (
	"context"
	"regexp"
	"strings"

	"vrcomandaapi/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	List(ctx context.Context, filter ListProdutosFilter) ([]Produto, error)
}

var findAllProdutosCodigoBarrasFn = database.FindAll[ProdutoCodigoBarras]
var findAllProdutosFn = database.FindAll[Produto]

type mongoRepository struct {
	getDatabase                func() *mongo.Database
	invalidateConnection       func()
	produtosCollectionName     string
	codigoBarrasCollectionName string
}

func NewMongoRepository(
	getDatabase func() *mongo.Database,
	invalidateConnection func(),
	produtosCollectionName string,
	codigoBarrasCollectionName string,
) Repository {
	return &mongoRepository{
		getDatabase:                getDatabase,
		invalidateConnection:       invalidateConnection,
		produtosCollectionName:     produtosCollectionName,
		codigoBarrasCollectionName: codigoBarrasCollectionName,
	}
}

func (r *mongoRepository) collection(collectionName string) *mongo.Collection {
	db := r.getDatabase()
	if db == nil {
		return nil
	}

	return db.Collection(collectionName)
}

func (r *mongoRepository) List(ctx context.Context, filter ListProdutosFilter) ([]Produto, error) {
	produtosQuery := bson.M{}
	if filter.IDLoja != 0 {
		produtosQuery["idLoja"] = filter.IDLoja
	}
	if filter.DescricaoCompleta != "" {
		pattern := ".*" + regexp.QuoteMeta(filter.DescricaoCompleta) + ".*"
		produtosQuery["descricaocompleta"] = bson.M{"$regex": pattern, "$options": "i"}
	}
	if filter.DescricaoCupom != "" {
		pattern := ".*" + regexp.QuoteMeta(filter.DescricaoCupom) + ".*"
		produtosQuery["descricaocupom"] = bson.M{"$regex": pattern, "$options": "i"}
	}

	codigoBarras := strings.TrimSpace(filter.CodigoBarras)
	idsFiltroCodigoBarras := []int{}
	if codigoBarras != "" {
		codigoBarrasCollection := r.collection(r.codigoBarrasCollectionName)

		codigos, err := findAllProdutosCodigoBarrasFn(
			ctx,
			codigoBarrasCollection,
			bson.M{"$or": []bson.M{{"codigobarras": codigoBarras}, {"codigoBarras": codigoBarras}}},
			nil,
			r.invalidateConnection,
		)
		if err != nil {
			return nil, err
		}
		if len(codigos) == 0 {
			return []Produto{}, nil
		}

		ids := make([]int, 0, len(codigos))
		seen := make(map[int]struct{}, len(codigos))
		for _, item := range codigos {
			if _, exists := seen[item.IDProduto]; exists {
				continue
			}
			seen[item.IDProduto] = struct{}{}
			ids = append(ids, item.IDProduto)
		}
		if len(ids) == 0 {
			return []Produto{}, nil
		}

		idsFiltroCodigoBarras = ids
		produtosQuery["idProduto"] = bson.M{"$in": ids}
	}

	findOptions := options.Find().SetSort(bson.M{"descricaocompleta": 1})
	produtosCollection := r.collection(r.produtosCollectionName)
	produtos, err := findAllProdutosFn(ctx, produtosCollection, produtosQuery, findOptions, r.invalidateConnection)
	if err != nil {
		return nil, err
	}
	if len(produtos) == 0 {
		return []Produto{}, nil
	}

	idsProdutos := idsFiltroCodigoBarras
	if len(idsProdutos) == 0 {
		idsProdutos = make([]int, 0, len(produtos))
		seen := make(map[int]struct{}, len(produtos))
		for _, item := range produtos {
			if _, exists := seen[item.IDProduto]; exists {
				continue
			}
			seen[item.IDProduto] = struct{}{}
			idsProdutos = append(idsProdutos, item.IDProduto)
		}
	}

	codigoBarrasCollection := r.collection(r.codigoBarrasCollectionName)
	codigosByProduto, err := findAllProdutosCodigoBarrasFn(
		ctx,
		codigoBarrasCollection,
		bson.M{"idProduto": bson.M{"$in": idsProdutos}},
		nil,
		r.invalidateConnection,
	)
	if err != nil {
		return nil, err
	}

	porProduto := make(map[int][]ProdutoCodigoBarras, len(idsProdutos))
	for _, item := range codigosByProduto {
		porProduto[item.IDProduto] = append(porProduto[item.IDProduto], item)
	}

	for i := range produtos {
		produtos[i].CodigosBarras = porProduto[produtos[i].IDProduto]
	}

	return produtos, nil
}
