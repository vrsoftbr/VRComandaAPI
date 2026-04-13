package mesa

type ListMesasRequest struct {
	IDLoja int   `form:"idLoja"`
	Mesa   int   `form:"mesa"`
	Ativo  *bool `form:"ativo"`
}

type MesaResponse struct {
	ID        string `json:"_id"`
	IDLoja    int    `json:"idLoja"`
	Mesa      int    `json:"mesa"`
	Descricao string `json:"descricao"`
	Ativo     bool   `json:"ativo"`
}

type ListMesasFilter struct {
	IDLoja int
	Mesa   int
	Ativo  *bool
}
