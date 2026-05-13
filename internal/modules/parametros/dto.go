package parametros

type ListParametrosRequest struct {
	IDLoja int `form:"idLoja"`
}

type ParametroResponse struct {
	IDParametro int    `json:"idParametro"`
	Descricao   string `json:"descricao"`
	IDLoja      int    `json:"idLoja"`
	Valor       string `json:"valor"`
}
