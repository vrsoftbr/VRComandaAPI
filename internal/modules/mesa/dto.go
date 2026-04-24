package mesa

type ListMesasRequest struct {
	IDLoja int   `form:"id_loja"`
	Mesa   int   `form:"mesa"`
	Mesas  []int `form:"mesas"`
	Ativo  *bool `form:"ativo"`
}

type MesaResponse struct {
	ID        string `json:"_id"`
	IDLoja    int    `json:"id_loja"`
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
