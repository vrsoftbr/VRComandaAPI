package parametros

import (
	"context"
	"errors"
	"testing"
)

type repositoryStub struct {
	listByLojaAndParametrosFn func(ctx context.Context, idLoja int, ids []int) ([]Parametro, error)
}

func (s repositoryStub) ListByLojaAndParametros(ctx context.Context, idLoja int, ids []int) ([]Parametro, error) {
	return s.listByLojaAndParametrosFn(ctx, idLoja, ids)
}

func TestServiceListRequiresIDLoja(t *testing.T) {
	svc := NewService()

	result, err := svc.List(context.Background(), ListParametrosRequest{})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestServiceListReturnsDefaultParametros(t *testing.T) {
	svc := NewServiceWithRepository(repositoryStub{listByLojaAndParametrosFn: func(_ context.Context, idLoja int, ids []int) ([]Parametro, error) {
		if idLoja != 12 {
			t.Fatalf("unexpected idLoja: %d", idLoja)
		}
		if len(ids) != len(defaultParametros) {
			t.Fatalf("expected %d ids, got %d", len(defaultParametros), len(ids))
		}

		return []Parametro{
			{IDParametro: 14, IDLoja: 12, Valor: "2"},
			{IDParametro: 19, IDLoja: 12, Valor: "5"},
		}, nil
	}})

	result, err := svc.List(context.Background(), ListParametrosRequest{IDLoja: 12})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(defaultParametros) {
		t.Fatalf("expected %d parametros, got %d", len(defaultParametros), len(result))
	}
	if result[0].IDParametro != 13 || result[0].IDLoja != 12 || result[0].Valor != "1" {
		t.Fatalf("unexpected first parametro: %+v", result[0])
	}
	if result[1].IDParametro != 14 || result[1].Valor != "2" {
		t.Fatalf("expected mongo value for parametro 14, got %+v", result[1])
	}
	if result[6].IDParametro != 19 || result[6].Valor != "5" {
		t.Fatalf("expected mongo value for parametro 19, got %+v", result[6])
	}
	if result[7].IDParametro != 227 || result[7].Valor != "1" {
		t.Fatalf("unexpected parametro 227: %+v", result[7])
	}
	if result[len(result)-1].IDParametro != 229 {
		t.Fatalf("unexpected last parametro: %+v", result[len(result)-1])
	}
}

func TestServiceListPropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repo failure")
	svc := NewServiceWithRepository(repositoryStub{listByLojaAndParametrosFn: func(_ context.Context, _ int, _ []int) ([]Parametro, error) {
		return nil, expectedErr
	}})

	result, err := svc.List(context.Background(), ListParametrosRequest{IDLoja: 1})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Fatalf("expected nil result on error, got %+v", result)
	}
}
