package loja

import "context"

type Service interface {
	List(ctx context.Context) ([]LojaResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context) ([]LojaResponse, error) {
	lojas, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]LojaResponse, 0, len(lojas))
	for _, l := range lojas {
		response = append(response, LojaResponse{
			ID:           l.ID,
			Descricao:    l.Descricao,
			RazaoSocial:  l.RazaoSocial,
			NomeFantasia: l.NomeFantasia,
			CodigoPais:   l.CodigoPais,
			Primaria:     l.Moedas.Primaria,
		})
	}

	return response, nil
}
