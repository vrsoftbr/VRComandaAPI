package produto

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	List(ctx context.Context, req ListProdutosRequest) (interface{}, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, req ListProdutosRequest) (interface{}, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	filter := ListProdutosFilter{
		IDLoja:            req.IDLoja,
		CodigoBarras:      req.CodigoBarras,
		DescricaoCompleta: req.DescricaoCompleta,
		DescricaoCupom:    req.DescricaoCupom,
		Page:              req.Page,
		Limit:             req.Limit,
	}

	result, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := make([]ProdutoResponse, 0, len(result.Items))
	for _, p := range result.Items {
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

	// Calcular total de páginas
	pages := int64(1)
	if result.Total > 0 {
		pages = (result.Total + int64(req.Limit) - 1) / int64(req.Limit)
	}

	return ProdutosPaginatedResponse{
		Items: response,
		Page:  req.Page,
		Limit: req.Limit,
		Total: result.Total,
		Pages: pages,
	}, nil
}

func decimalToString(value primitive.Decimal128) string {
	return value.String()
}
