package produto

type ListProdutosRequest struct {
	IDLoja            int    `form:"idLoja"`
	CodigoBarras      string `form:"codigoBarras"`
	DescricaoCompleta string `form:"descricaocompleta"`
	DescricaoCupom    string `form:"descricaocupom"`
	Page              int    `form:"page" binding:"omitempty,gt=0"`
	Limit             int    `form:"limit" binding:"omitempty,gt=0,lte=100"`
}

type ProdutoCodigoBarrasResponse struct {
	ID                  int32  `json:"_id"`
	IDProduto           int    `json:"idProduto"`
	CodigoBarras        string `json:"codigobarras"`
	Embalagem           string `json:"embalagem"`
	QuantidadeEmbalagem int    `json:"quantidadeembalagem"`
	Classe              string `json:"_class,omitempty"`
}

type ProdutoResponse struct {
	ID                    string `json:"_id"`
	IDProduto             int    `json:"idProduto"`
	IDLoja                int    `json:"idLoja"`
	DescricaoCompleta     string `json:"descricaocompleta"`
	DescricaoCupom        string `json:"descricaocupom"`
	PrecoVenda            string `json:"precovenda"`
	PrecoEspecial         string `json:"precoespecial"`
	PrecoEstrategico      string `json:"precoestrategico"`
	PermiteMultiplicacao  bool   `json:"permitemultiplicacao"`
	VendaControlada       bool   `json:"vendacontrolada"`
	QuantidadeParcela     int    `json:"quantidadeparcela"`
	DescontoMaximo        string `json:"descontomaximo"`
	ValidaPeso            bool   `json:"validapeso"`
	Pesavel               bool   `json:"pesavel"`
	IDProdutoVasilhame    int    `json:"idProdutoVasilhame"`
	NCM                   string `json:"ncm"`
	CEST                  string `json:"cest"`
	OrigemMercadoriaSaida int    `json:"origemmercadoriasaida"`
	IDCenarioFiscal       int    `json:"idcenariofiscal"`
	CodigoBarras          string `json:"codigobarras"`
	Embalagem             string `json:"embalagem"`
	QuantidadeEmbalagem   int    `json:"quantidadeembalagem"`
	Classe                string `json:"_class,omitempty"`
}

type ListProdutosFilter struct {
	IDLoja            int
	CodigoBarras      string
	DescricaoCompleta string
	DescricaoCupom    string
	Page              int
	Limit             int
}

type ProdutosPaginatedResponse struct {
	Items []ProdutoResponse `json:"items"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Total int64             `json:"total"`
	Pages int64             `json:"pages"`
}
