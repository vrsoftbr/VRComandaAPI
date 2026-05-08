package mesa

type ListMesasRequest struct {
	IDLoja int   `form:"idLoja"`
	Mesa   int   `form:"mesa"`
	Mesas  []int `form:"mesas"`
	Ativo  *bool `form:"ativo"`
}

type MesaResponse struct {
	IDLoja    int    `json:"idLoja"`
	Mesa      int    `json:"mesa"`
	Descricao string `json:"descricao"`
	Ativo     bool   `json:"ativo"`
}

type ListMesasFilter struct {
	IDLoja int
	Mesa   int
	Mesas  []int
	Ativo  *bool
}
