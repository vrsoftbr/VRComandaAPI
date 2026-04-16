package comanda

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type repositoryStub struct {
	listFn func(ctx context.Context, filter ListComandasFilter) ([]Comanda, error)
}

func (s repositoryStub) List(ctx context.Context, filter ListComandasFilter) ([]Comanda, error) {
	return s.listFn(ctx, filter)
}

func TestServiceListBuildsFilterAndMapsResponse(t *testing.T) {
	ativo := true
	expectedID := primitive.NewObjectID()

	called := 0
	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListComandasFilter) ([]Comanda, error) {
			called++
			return []Comanda{{
				ID:                  expectedID,
				IDLoja:              20,
				Comanda:             101,
				NumeroIdentificacao: "1",
				Observacao:          "obs",
				Ativo:               true,
			}}, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListComandasRequest{
		IDLoja:              20,
		Comanda:             101,
		NumeroIdentificacao: "1",
		Ativo:               &ativo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Fatalf("repo.List calls = %d", called)
	}
	if len(result) != 1 {
		t.Fatalf("expected one item, got %d", len(result))
	}
}

func TestServiceListDoesNotSetAtivoWhenNil(t *testing.T) {
	var capturedFilter ListComandasFilter
	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListComandasFilter) ([]Comanda, error) {
			capturedFilter = filter
			return []Comanda{}, nil
		},
	}

	svc := NewService(repo)
	_, err := svc.List(context.Background(), ListComandasRequest{IDLoja: 1, Comanda: 2, NumeroIdentificacao: "N"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedFilter.Ativo != nil {
		t.Fatalf("expected nil ativo in filter, got %+v", capturedFilter.Ativo)
	}
}

func TestServiceListPropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repo failure")
	repo := repositoryStub{
		listFn: func(_ context.Context, _ ListComandasFilter) ([]Comanda, error) {
			return nil, expectedErr
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListComandasRequest{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Fatalf("expected nil result on error, got %+v", result)
	}
}
