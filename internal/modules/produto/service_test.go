package produto

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type repositoryStub struct {
	listFn func(ctx context.Context, filter ListProdutosFilter) ([]Produto, error)
}

func (s repositoryStub) List(ctx context.Context, filter ListProdutosFilter) ([]Produto, error) {
	return s.listFn(ctx, filter)
}

func TestServiceListBuildsFilterAndMapsResponse(t *testing.T) {
	precoVenda, _ := primitive.ParseDecimal128("3.49")
	precoEspecial, _ := primitive.ParseDecimal128("0")
	precoEstrategico, _ := primitive.ParseDecimal128("0")
	descontoMaximo, _ := primitive.ParseDecimal128("0.00")
	var capturedFilter ListProdutosFilter

	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListProdutosFilter) ([]Produto, error) {
			capturedFilter = filter
			return []Produto{{
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
			}}, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListProdutosRequest{
		IDLoja:            5,
		CodigoBarras:      "123",
		DescricaoCompleta: "REPOLHO",
		DescricaoCupom:    "KG",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedFilter.IDLoja != 5 || capturedFilter.CodigoBarras != "123" || capturedFilter.DescricaoCompleta != "REPOLHO" || capturedFilter.DescricaoCupom != "KG" {
		t.Fatalf("unexpected filter passed to repository: %+v", capturedFilter)
	}
	if len(result) != 1 {
		t.Fatalf("expected one item, got %d", len(result))
	}
	if result[0].IDProduto != 3 || result[0].IDLoja != 5 {
		t.Fatalf("unexpected mapped response: %+v", result[0])
	}
	if result[0].PrecoVenda != "3.49" || result[0].DescontoMaximo != "0.00" || !result[0].PermiteMultiplicacao {
		t.Fatalf("unexpected full field mapping: %+v", result[0])
	}
	if result[0].CodigoBarras != "1" || result[0].Embalagem != "KG" || result[0].QuantidadeEmbalagem != 1 {
		t.Fatalf("unexpected barcode mapping: %+v", result[0])
	}
}

func TestServiceListPropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repo failure")
	repo := repositoryStub{
		listFn: func(_ context.Context, _ ListProdutosFilter) ([]Produto, error) {
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
