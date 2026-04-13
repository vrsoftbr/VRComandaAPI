package atendente

type ListAtendentesRequest struct {
	IDLoja int    `form:"idLoja"`
	Codigo string `form:"codigo"`
	Nome   string `form:"nome"`
	Ativo  *bool  `form:"ativo"`
}

type AtendenteResponse struct {
	ID     string `json:"_id"`
	IDLoja int    `json:"idLoja"`
	Codigo string `json:"codigo"`
	Nome   string `json:"nome"`
	Senha  string `json:"senha"`
	Ativo  bool   `json:"ativo"`
}

type ListAtendentesFilter struct {
	IDLoja int
	Codigo string
	Nome   string
	Ativo  *bool
}
