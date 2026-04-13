package mesa

import (
	"context"
)

type Service interface {
	List(ctx context.Context, req ListMesasRequest) ([]MesaResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, req ListMesasRequest) ([]MesaResponse, error) {
	filter := ListMesasFilter{
		IDLoja: req.IDLoja,
		Mesa:   req.Mesa,
	}

	if req.Ativo != nil {
		filter.Ativo = req.Ativo
	}

	models, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := make([]MesaResponse, 0, len(models))
	for _, m := range models {
		response = append(response, MesaResponse{
			ID:        m.ID.Hex(),
			IDLoja:    m.IDLoja,
			Mesa:      m.Mesa,
			Descricao: m.Descricao,
			Ativo:     m.Ativo,
		})
	}

	return response, nil
}
