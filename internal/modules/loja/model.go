package loja

type Moedas struct {
	Primaria int `bson:"primaria"`
}

type Loja struct {
	ID           int    `bson:"_id"`
	Descricao    string `bson:"descricao"`
	RazaoSocial  string `bson:"razaosocial"`
	NomeFantasia string `bson:"nomefantasia"`
	CodigoPais   int    `bson:"codigopais"`
	Moedas       Moedas `bson:"moedas"`
}
