package produto

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	List(ctx context.Context, req ListProdutosRequest) ([]ProdutoResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, req ListProdutosRequest) ([]ProdutoResponse, error) {
	filter := ListProdutosFilter{
		IDLoja:            req.IDLoja,
		CodigoBarras:      req.CodigoBarras,
		DescricaoCompleta: req.DescricaoCompleta,
		DescricaoCupom:    req.DescricaoCupom,
	}

	models, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := make([]ProdutoResponse, 0, len(models))
	for _, p := range models {
		var codigoBarras string
		var embalagem string
		var quantidadeEmbalagem int
		if len(p.CodigosBarras) > 0 {
			cb := p.CodigosBarras[0]
			codigoBarras = cb.CodigoBarras
			if codigoBarras == "" {
				codigoBarras = cb.CodigoBarrasLegacy
			}
			embalagem = cb.Embalagem
			quantidadeEmbalagem = cb.QuantidadeEmbalagem
		}

		response = append(response, ProdutoResponse{
			IDProduto:             p.IDProduto,
			IDLoja:                p.IDLoja,
			DescricaoCompleta:     p.DescricaoCompleta,
			DescricaoCupom:        p.DescricaoCupom,
			PrecoVenda:            decimalToString(p.PrecoVenda),
			PrecoEspecial:         decimalToString(p.PrecoEspecial),
			PrecoEstrategico:      decimalToString(p.PrecoEstrategico),
			PermiteMultiplicacao:  p.PermiteMultiplicacao,
			VendaControlada:       p.VendaControlada,
			QuantidadeParcela:     p.QuantidadeParcela,
			DescontoMaximo:        decimalToString(p.DescontoMaximo),
			ValidaPeso:            p.ValidaPeso,
			Pesavel:               p.Pesavel,
			IDProdutoVasilhame:    p.IDProdutoVasilhame,
			NCM:                   p.NCM,
			CEST:                  p.CEST,
			OrigemMercadoriaSaida: p.OrigemMercadoriaSaida,
			IDCenarioFiscal:       p.IDCenarioFiscal,
			CodigoBarras:          codigoBarras,
			Embalagem:             embalagem,
			QuantidadeEmbalagem:   quantidadeEmbalagem,
		})
	}

	return response, nil
}

func decimalToString(value primitive.Decimal128) string {
	return value.String()
}
