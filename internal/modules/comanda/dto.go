package comanda

type ListComandasRequest struct {
	IDLoja              int    `form:"idLoja"`
	Comanda             int    `form:"comanda"`
	NumeroIdentificacao string `form:"numeroIdentificacao"`
	Ativo               *bool  `form:"ativo"`
}

type ComandaResponse struct {
	ID                  string `json:"_id"`
	IDLoja              int    `json:"idLoja"`
	Comanda             int    `json:"comanda"`
	NumeroIdentificacao string `json:"numeroIdentificacao"`
	Observacao          string `json:"observacao"`
	Ativo               bool   `json:"ativo"`
}

type ListComandasFilter struct {
	IDLoja              int
	Comanda             int
	NumeroIdentificacao string
	Ativo               *bool
}
