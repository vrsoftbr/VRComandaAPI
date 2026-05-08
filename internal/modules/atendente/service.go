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
		IDLoja:      req.IDLoja,
		IDAtendente: req.IDAtendente,
		Nome:        req.Nome,
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
			IDLoja:      m.IDLoja,
			IDAtendente: m.IDAtendente,
			Nome:        m.Nome,
			Senha:       m.Senha,
			Ativo:       m.Ativo,
		})
	}

	return response, nil
}
