package atendente

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type repositoryStub struct {
	listFn func(ctx context.Context, filter ListAtendentesFilter) ([]Atendente, error)
}

func (s repositoryStub) List(ctx context.Context, filter ListAtendentesFilter) ([]Atendente, error) {
	return s.listFn(ctx, filter)
}

func TestServiceListBuildsFilterAndMapsResponse(t *testing.T) {
	ativo := true
	expectedID := primitive.NewObjectID()

	called := 0
	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListAtendentesFilter) ([]Atendente, error) {
			called++
			return []Atendente{{
				ID:     expectedID,
				IDLoja: 20,
				Codigo: "01",
				Nome:   "Maria",
				Senha:  "123",
				Ativo:  true,
			}}, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListAtendentesRequest{
		IDLoja: 20,
		Codigo: "01",
		Nome:   "Mar",
		Ativo:  &ativo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected one item, got %d", len(result))
	}
}

func TestServiceListDoesNotSetAtivoWhenNil(t *testing.T) {
	var capturedFilter ListAtendentesFilter
	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListAtendentesFilter) ([]Atendente, error) {
			capturedFilter = filter
			return []Atendente{}, nil
		},
	}

	svc := NewService(repo)
	_, err := svc.List(context.Background(), ListAtendentesRequest{IDLoja: 1, Codigo: "C", Nome: "N"})
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
		listFn: func(_ context.Context, _ ListAtendentesFilter) ([]Atendente, error) {
			return nil, expectedErr
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListAtendentesRequest{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Fatalf("expected nil result on error, got %+v", result)
	}
}
