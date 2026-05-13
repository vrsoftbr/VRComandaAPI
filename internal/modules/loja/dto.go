package loja

type LojaResponse struct {
	ID           int    `json:"_id"`
	Descricao    string `json:"descricao"`
	RazaoSocial  string `json:"razaosocial"`
	NomeFantasia string `json:"nomefantasia"`
	CodigoPais   int    `json:"codigopais"`
	Primaria     int    `json:"primaria"`
}
