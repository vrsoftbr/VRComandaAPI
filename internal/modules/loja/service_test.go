package loja

import (
	"context"
	"errors"
	"testing"
)

type repositoryStub struct {
	listFn func(ctx context.Context) ([]Loja, error)
}

func (s repositoryStub) List(ctx context.Context) ([]Loja, error) {
	return s.listFn(ctx)
}

func TestServiceListPropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repo failure")
	svc := NewService(repositoryStub{listFn: func(_ context.Context) ([]Loja, error) {
		return nil, expectedErr
	}})

	result, err := svc.List(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Fatalf("expected nil result on error, got %+v", result)
	}
}

func TestServiceListMapsResponseCorrectly(t *testing.T) {
	svc := NewService(repositoryStub{listFn: func(_ context.Context) ([]Loja, error) {
		return []Loja{
			{
				ID:           1,
				Descricao:    "Loja Um",
				RazaoSocial:  "Razao Social",
				NomeFantasia: "Fantasia",
				CodigoPais:   1058,
				Moedas:       Moedas{Primaria: 5},
			},
		}, nil
	}})

	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected one item, got %d", len(result))
	}
	r := result[0]
	if r.ID != 1 || r.Descricao != "Loja Um" || r.RazaoSocial != "Razao Social" || r.NomeFantasia != "Fantasia" || r.CodigoPais != 1058 || r.Primaria != 5 {
		t.Fatalf("unexpected result: %+v", r)
	}
}

func TestServiceListReturnsEmpty(t *testing.T) {
	svc := NewService(repositoryStub{listFn: func(_ context.Context) ([]Loja, error) {
		return []Loja{}, nil
	}})

	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d items", len(result))
	}
}
