package comanda

import (
	"context"
)

type Service interface {
	List(ctx context.Context, req ListComandasRequest) ([]ComandaResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, req ListComandasRequest) ([]ComandaResponse, error) {
	filter := ListComandasFilter{
		IDLoja:              req.IDLoja,
		Comanda:             req.Comanda,
		NumeroIdentificacao: req.NumeroIdentificacao,
	}

	if req.Ativo != nil {
		filter.Ativo = req.Ativo
	}

	models, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := make([]ComandaResponse, 0, len(models))
	for _, m := range models {
		response = append(response, ComandaResponse{
			ID:                  m.ID.Hex(),
			IDLoja:              m.IDLoja,
			Comanda:             m.Comanda,
			NumeroIdentificacao: m.NumeroIdentificacao,
			Observacao:          m.Observacao,
			Ativo:               m.Ativo,
		})
	}

	return response, nil
}
