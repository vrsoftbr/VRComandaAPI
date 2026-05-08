package comanda

type Comanda struct {
	IDLoja              int    `bson:"idLoja" json:"idLoja"`
	Comanda             int    `bson:"comanda" json:"comanda"`
	NumeroIdentificacao string `bson:"numeroIdentificacao" json:"numeroIdentificacao"`
	Observacao          string `bson:"observacao" json:"observacao"`
	Ativo               bool   `bson:"ativo" json:"ativo"`
}
