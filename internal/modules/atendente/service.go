package atendente

import (
	"context"
)

type Service interface {
	List(ctx context.Context, req ListAtendentesRequest) ([]AtendenteResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, req ListAtendentesRequest) ([]AtendenteResponse, error) {
	filter := ListAtendentesFilter{
		IDLoja: req.IDLoja,
		Codigo: req.Codigo,
		Nome:   req.Nome,
	}

	if req.Ativo != nil {
		filter.Ativo = req.Ativo
	}

	models, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := make([]AtendenteResponse, 0, len(models))
	for _, m := range models {
		response = append(response, AtendenteResponse{
			ID:     m.ID.Hex(),
			IDLoja: m.IDLoja,
			Codigo: m.Codigo,
			Nome:   m.Nome,
			Senha:  m.Senha,
			Ativo:  m.Ativo,
		})
	}

	return response, nil
}
