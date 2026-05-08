package atendente

type ListAtendentesRequest struct {
	IDLoja      int    `form:"idLoja"`
	IDAtendente string `form:"idAtendente"`
	Nome        string `form:"nome"`
	Ativo       *bool  `form:"ativo"`
}

type AtendenteResponse struct {
	IDLoja      int    `json:"idLoja"`
	IDAtendente string `json:"idAtendente"`
	Nome        string `json:"nome"`
	Senha       string `json:"senha"`
	Ativo       bool   `json:"ativo"`
}

type ListAtendentesFilter struct {
	IDLoja      int
	IDAtendente string
	Nome        string
	Ativo       *bool
}
