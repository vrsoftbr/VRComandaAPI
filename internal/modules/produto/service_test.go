package produto

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type repositoryStub struct {
	listFn func(ctx context.Context, filter ListProdutosFilter) (*RepositoryResult, error)
}

func (s repositoryStub) List(ctx context.Context, filter ListProdutosFilter) (*RepositoryResult, error) {
	return s.listFn(ctx, filter)
}

func TestServiceListBuildsFilterAndMapsResponse(t *testing.T) {
	precoVenda, _ := primitive.ParseDecimal128("3.49")
	precoEspecial, _ := primitive.ParseDecimal128("0")
	precoEstrategico, _ := primitive.ParseDecimal128("0")
	descontoMaximo, _ := primitive.ParseDecimal128("0.00")
	var capturedFilter ListProdutosFilter

	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListProdutosFilter) (*RepositoryResult, error) {
			capturedFilter = filter
			return &RepositoryResult{
				Items: []Produto{{
					IDProduto:             3,
					IDLoja:                5,
					DescricaoCompleta:     "REPOLHO ROXO KG",
					DescricaoCupom:        "REPOLHO ROXO KG",
					PrecoVenda:            precoVenda,
					PrecoEspecial:         precoEspecial,
					PrecoEstrategico:      precoEstrategico,
					PermiteMultiplicacao:  true,
					VendaControlada:       false,
					QuantidadeParcela:     0,
					DescontoMaximo:        descontoMaximo,
					ValidaPeso:            false,
					Pesavel:               true,
					IDProdutoVasilhame:    2,
					NCM:                   "08094000",
					CEST:                  "0300700",
					OrigemMercadoriaSaida: 0,
					IDCenarioFiscal:       335,
					CodigosBarras: []ProdutoCodigoBarras{{
						IDProduto:           3,
						CodigoBarras:        "1",
						Embalagem:           "KG",
						QuantidadeEmbalagem: 1,
					}},
				}},
				Total: 1,
			}, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListProdutosRequest{
		IDLoja:            5,
		CodigoBarras:      "123",
		DescricaoCompleta: "REPOLHO",
		DescricaoCupom:    "KG",
		Page:              1,
		Limit:             20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedFilter.IDLoja != 5 || capturedFilter.CodigoBarras != "123" || capturedFilter.DescricaoCompleta != "REPOLHO" || capturedFilter.DescricaoCupom != "KG" {
		t.Fatalf("unexpected filter passed to repository: %+v", capturedFilter)
	}

	paginatedResult, ok := result.(ProdutosPaginatedResponse)
	if !ok {
		t.Fatalf("expected ProdutosPaginatedResponse, got %T", result)
	}
	if len(paginatedResult.Items) != 1 {
		t.Fatalf("expected one item, got %d", len(paginatedResult.Items))
	}
	if paginatedResult.Items[0].IDProduto != 3 || paginatedResult.Items[0].IDLoja != 5 {
		t.Fatalf("unexpected mapped response: %+v", paginatedResult.Items[0])
	}
	if paginatedResult.Items[0].PrecoVenda != "3.49" || paginatedResult.Items[0].DescontoMaximo != "0.00" || !paginatedResult.Items[0].PermiteMultiplicacao {
		t.Fatalf("unexpected full field mapping: %+v", paginatedResult.Items[0])
	}
	if paginatedResult.Items[0].CodigoBarras != "1" || paginatedResult.Items[0].Embalagem != "KG" || paginatedResult.Items[0].QuantidadeEmbalagem != 1 {
		t.Fatalf("unexpected barcode mapping: %+v", paginatedResult.Items[0])
	}
	if paginatedResult.Total != 1 || paginatedResult.Pages != 1 {
		t.Fatalf("unexpected pagination: %+v", paginatedResult)
	}
}

func TestServiceListPropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repo failure")
	repo := repositoryStub{
		listFn: func(_ context.Context, _ ListProdutosFilter) (*RepositoryResult, error) {
			return nil, expectedErr
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListProdutosRequest{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Fatalf("expected nil result on error, got %+v", result)
	}
}

func TestServiceListFallsBackToCodigoBarrasLegacy(t *testing.T) {
	repo := repositoryStub{
		listFn: func(_ context.Context, _ ListProdutosFilter) (*RepositoryResult, error) {
			return &RepositoryResult{
				Items: []Produto{{
					IDProduto: 1,
					IDLoja:    1,
					CodigosBarras: []ProdutoCodigoBarras{{
						IDProduto:           1,
						CodigoBarras:        "", // empty → should use legacy
						CodigoBarrasLegacy:  "LEG123",
						Embalagem:           "UN",
						QuantidadeEmbalagem: 1,
					}},
				}},
				Total: 1,
			}, nil
		},
	}

	svc := NewService(repo)
	raw, err := svc.List(context.Background(), ListProdutosRequest{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := raw.(ProdutosPaginatedResponse)
	if len(result.Items) != 1 || result.Items[0].CodigoBarras != "LEG123" {
		t.Fatalf("expected CodigoBarras=LEG123, got %+v", result.Items)
	}
}

func TestServiceListWithNoCodigosBarras(t *testing.T) {
	repo := repositoryStub{
		listFn: func(_ context.Context, _ ListProdutosFilter) (*RepositoryResult, error) {
			return &RepositoryResult{
				Items: []Produto{{
					IDProduto:     2,
					IDLoja:        1,
					CodigosBarras: []ProdutoCodigoBarras{},
				}},
				Total: 0,
			}, nil
		},
	}

	svc := NewService(repo)
	raw, err := svc.List(context.Background(), ListProdutosRequest{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := raw.(ProdutosPaginatedResponse)
	if len(result.Items) != 1 || result.Items[0].CodigoBarras != "" {
		t.Fatalf("expected empty CodigoBarras, got %+v", result.Items)
	}
}
