package mesa

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type repositoryStub struct {
	listFn func(ctx context.Context, filter ListMesasFilter) ([]Mesa, error)
}

func (s repositoryStub) List(ctx context.Context, filter ListMesasFilter) ([]Mesa, error) {
	return s.listFn(ctx, filter)
}

func TestServiceListBuildsFilterAndMapsResponse(t *testing.T) {
	ativo := true
	expectedID := primitive.NewObjectID()

	called := 0
	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListMesasFilter) ([]Mesa, error) {
			called++
			return []Mesa{{
				ID:        expectedID,
				IDLoja:    20,
				Mesa:      8,
				Descricao: "mesa 8",
				Ativo:     true,
			}}, nil
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListMesasRequest{
		IDLoja: 20,
		Mesa:   8,
		Mesas:  []int{8, 9},
		Ativo:  &ativo,
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
	var capturedFilter ListMesasFilter
	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListMesasFilter) ([]Mesa, error) {
			capturedFilter = filter
			return []Mesa{}, nil
		},
	}

	svc := NewService(repo)
	_, err := svc.List(context.Background(), ListMesasRequest{IDLoja: 1, Mesa: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedFilter.Ativo != nil {
		t.Fatalf("expected nil ativo in filter, got %+v", capturedFilter.Ativo)
	}
	if len(capturedFilter.Mesas) != 0 {
		t.Fatalf("expected empty mesas in filter, got %+v", capturedFilter.Mesas)
	}
}

func TestServiceListPassesBatchMesas(t *testing.T) {
	var capturedFilter ListMesasFilter
	repo := repositoryStub{
		listFn: func(_ context.Context, filter ListMesasFilter) ([]Mesa, error) {
			capturedFilter = filter
			return []Mesa{}, nil
		},
	}

	svc := NewService(repo)
	_, err := svc.List(context.Background(), ListMesasRequest{Mesas: []int{4, 6}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(capturedFilter.Mesas) != 2 || capturedFilter.Mesas[0] != 4 || capturedFilter.Mesas[1] != 6 {
		t.Fatalf("unexpected mesas filter: %+v", capturedFilter.Mesas)
	}
}

func TestServiceListPropagatesRepositoryError(t *testing.T) {
	expectedErr := errors.New("repo failure")
	repo := repositoryStub{
		listFn: func(_ context.Context, _ ListMesasFilter) ([]Mesa, error) {
			return nil, expectedErr
		},
	}

	svc := NewService(repo)
	result, err := svc.List(context.Background(), ListMesasRequest{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Fatalf("expected nil result on error, got %+v", result)
	}
}
