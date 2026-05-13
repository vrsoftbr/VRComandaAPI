package parametros

import (
	"context"
	"fmt"
)

type Service interface {
	List(ctx context.Context, req ListParametrosRequest) ([]ParametroResponse, error)
}

type service struct {
	repo Repository
}

type parametroTemplate struct {
	idParametro int
	descricao   string
	valor       string
}

var defaultParametros = []parametroTemplate{
	{idParametro: 13, descricao: "Tipo de Etiqueta de Balanca", valor: "1"},
	{idParametro: 14, descricao: "Posicao Inicial do Codigo do Produto na Etiqueta de Balanca", valor: "2"},
	{idParametro: 15, descricao: "Quantidade de Digitos do Codigo do Produto na Etiqueta de Balanca", valor: "5"},
	{idParametro: 16, descricao: "Posicao Inicial do Preco do Produto na Etiqueta de Balanca", valor: "7"},
	{idParametro: 17, descricao: "Quantidade de Digitos do Preco do Produto na Etiqueta de Balanca", valor: "6"},
	{idParametro: 18, descricao: "Posicao Inicial do Peso do Produto na Etiqueta de Balanca", valor: "2"},
	{idParametro: 19, descricao: "Quantidade de Digitos do Peso do Produto na Etiqueta de Balanca", valor: "5"},
	{idParametro: 227, descricao: "Informa Atendente no Lançamento da Comanda", valor: "1"},
	{idParametro: 228, descricao: "Informa Mesa no Lancamento da Comanda", valor: "1"},
	{idParametro: 229, descricao: "Tipo de Identificacao da Comanda", valor: "1"},
}

func NewService() Service {
	return &service{}
}

func NewServiceWithRepository(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, req ListParametrosRequest) ([]ParametroResponse, error) {
	if req.IDLoja <= 0 {
		return nil, fmt.Errorf("%w: idLoja deve ser maior que zero", ErrInvalidRequest)
	}

	ids := make([]int, 0, len(defaultParametros))
	for _, item := range defaultParametros {
		ids = append(ids, item.idParametro)
	}

	valuesByParametro := map[int]string{}
	if s.repo != nil {
		models, err := s.repo.ListByLojaAndParametros(ctx, req.IDLoja, ids)
		if err != nil {
			return nil, err
		}

		for _, model := range models {
			valuesByParametro[model.IDParametro] = model.Valor
		}
	}

	response := make([]ParametroResponse, 0, len(defaultParametros))
	for _, item := range defaultParametros {
		valor := item.valor
		if v, ok := valuesByParametro[item.idParametro]; ok {
			valor = v
		}

		response = append(response, ParametroResponse{
			IDParametro: item.idParametro,
			Descricao:   item.descricao,
			IDLoja:      req.IDLoja,
			Valor:       valor,
		})
	}

	return response, nil
}
