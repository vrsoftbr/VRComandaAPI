package lancamento

import "time"

type ListItensRequest struct {
	IDComanda string `form:"id_comanda" binding:"required"`
}

type ItemComandaResponse struct {
	ID                   uint       `json:"id"`
	IDLancamentoComanda  uint       `json:"id_lancamentocomanda"`
	IDComanda            int        `json:"id_comanda"`
	Sequencia            int        `json:"sequencia"`
	IDProduto            int        `json:"id_produto"`
	CodigoBarras         string     `json:"codigobarras,omitempty"`
	Quantidade           float64    `json:"quantidade"`
	PrecoVenda           float64    `json:"precovenda"`
	Cancelado            bool       `json:"cancelado"`
	DataHoraCancelamento *time.Time `json:"dataHoraCancelamento,omitempty"`
	IDAtendente          int        `json:"id_atendente"`
	IDSituacao           int        `json:"id_situacao"`
}

type UpdateItemRequest struct {
	Quantidade *float64 `json:"quantidade"`
	Cancelado  *bool    `json:"cancelado"`
}

type CreateItemRequest struct {
	IDLancamentoComanda uint    `json:"id_lancamentocomanda" binding:"required"`
	Sequencia           int     `json:"sequencia" binding:"required"`
	IDProduto           int     `json:"id_produto" binding:"required"`
	CodigoBarras        string  `json:"codigobarras"`
	Quantidade          float64 `json:"quantidade" binding:"required"`
	PrecoVenda          float64 `json:"precovenda"`
	IDAtendente         int     `json:"id_atendente" binding:"required"`
	IDSituacao          int     `json:"id_situacao"`
}

type CreateItemsRequest struct {
	Itens []CreateItemRequest `json:"itens" binding:"required,min=1"`
}

type CreateLancamentoRequest struct {
	IDLoja      int    `json:"id_loja"`
	IDComanda   int    `json:"id_comanda"`
	IDMesa      *int   `json:"id_mesa"`
	IDAtendente int    `json:"id_atendente"`
	DataHora    string `json:"dataHora"`
	Observacao  string `json:"observacao"`
	Finalizado  *bool  `json:"finalizado"`
}

type ListLancamentosRequest struct {
	IDComanda   string `form:"id_comanda"`
	IDMesa      string `form:"id_mesa"`
	IDAtendente string `form:"id_atendente"`
	DataHora    string `form:"dataHora"`
	Finalizado  *bool  `form:"finalizado"`
}

type ListLancamentosFilter struct {
	IDComanda   *int
	IDMesa      *int
	IDAtendente *int
	DataHora    *time.Time
	Finalizado  *bool
}
