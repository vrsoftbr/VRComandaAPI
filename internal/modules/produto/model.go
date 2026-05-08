package produto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Produto struct {
	IDProduto             int                   `bson:"idProduto" json:"idProduto"`
	IDLoja                int                   `bson:"idLoja" json:"idLoja"`
	DescricaoCompleta     string                `bson:"descricaocompleta" json:"descricaocompleta"`
	DescricaoCupom        string                `bson:"descricaocupom" json:"descricaocupom"`
	PrecoVenda            primitive.Decimal128  `bson:"precovenda" json:"precovenda"`
	PrecoEspecial         primitive.Decimal128  `bson:"precoespecial" json:"precoespecial"`
	PrecoEstrategico      primitive.Decimal128  `bson:"precoestrategico" json:"precoestrategico"`
	PermiteMultiplicacao  bool                  `bson:"permitemultiplicacao" json:"permitemultiplicacao"`
	VendaControlada       bool                  `bson:"vendacontrolada" json:"vendacontrolada"`
	QuantidadeParcela     int                   `bson:"quantidadeparcela" json:"quantidadeparcela"`
	DescontoMaximo        primitive.Decimal128  `bson:"descontomaximo" json:"descontomaximo"`
	ValidaPeso            bool                  `bson:"validapeso" json:"validapeso"`
	Pesavel               bool                  `bson:"pesavel" json:"pesavel"`
	IDProdutoVasilhame    int                   `bson:"idProdutoVasilhame" json:"idProdutoVasilhame"`
	NCM                   string                `bson:"ncm" json:"ncm"`
	CEST                  string                `bson:"cest" json:"cest"`
	OrigemMercadoriaSaida int                   `bson:"origemmercadoriasaida" json:"origemmercadoriasaida"`
	IDCenarioFiscal       int                   `bson:"idcenariofiscal" json:"idcenariofiscal"`
	CodigosBarras         []ProdutoCodigoBarras `bson:"-" json:"produtoscodigobarras"`
}

type ProdutoCodigoBarras struct {
	IDProduto           int    `bson:"idProduto"`
	CodigoBarras        string `bson:"codigobarras"`
	CodigoBarrasLegacy  string `bson:"codigoBarras"`
	Embalagem           string `bson:"embalagem"`
	QuantidadeEmbalagem int    `bson:"quantidadeembalagem"`
}
