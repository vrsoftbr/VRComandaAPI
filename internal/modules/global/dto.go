package global

import (
	"vrcomandaapi/internal/modules/atendente"
	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/shared/models"
)

type LancamentoDetalhesDTO struct {
	IDLancamento uint                         `json:"id_lancamento"`
	IDLoja       int                          `json:"id_loja"`
	IDComanda    int                          `json:"id_comanda"`
	IDMesa       *int                         `json:"id_mesa,omitempty"`
	DataHora     string                       `json:"dataHora"`
	Observacao   string                       `json:"observacao,omitempty"`
	IDAtendente  int                          `json:"id_atendente"`
	Finalizado   bool                         `json:"finalizado"`
	Atendente    *atendente.AtendenteResponse `json:"atendente,omitempty"`
	Comanda      *comanda.ComandaResponse     `json:"comanda,omitempty"`
	Mesa         *mesa.MesaResponse           `json:"mesa,omitempty"`
	Itens        []LancamentoDetalhesItemDTO  `json:"itens,omitempty"`
}

type LancamentoDetalhesItemDTO struct {
	models.LancamentoComandaItem
	DescricaoProduto string `json:"descricaoProduto,omitempty"`
}

type GlobalFilterRequest struct {
	IDLoja     int   `form:"idLoja"`
	Finalizado *bool `form:"finalizado"`
}

type ListLancamentosDetalhesRequest struct {
	GlobalFilterRequest
}

// Usado pela catraca para liberar ou bloquear a comanda
type SituacaoComanda int

const (
	SituacaoComandaLiberada  SituacaoComanda = 1
	SituacaoComandaBloqueada SituacaoComanda = 2
)

type ComandaCatracaRequest struct {
	IDLoja                     int    `form:"idLoja"`
	NumeroIdentificacaoComanda string `form:"numeroIdentificacaoComanda"`
}

type ComandaCatracaResponse struct {
	IDLoja                     int             `json:"idLoja"`
	Comanda                    int             `json:"comanda"`
	NumeroIdentificacaoComanda string          `json:"numeroIdentificacaoComanda"`
	Situacao                   SituacaoComanda `json:"situacao"`
}

type GetLancamentoPDVRequest struct {
	NumeroComanda int `form:"numeroComanda"`
	IDLoja        int `form:"loja"`
}

type GetLancamentoItemPDVResponse struct {
	CodigoComanda        int                       `json:"codigoComanda"`
	TipoDocumentoCliente int                       `json:"tipoDocumentoCliente"`
	DocumentoCliente     string                    `json:"documentoCliente"`
	NomeCliente          string                    `json:"nomeCliente"`
	CodigoVendedor       int                       `json:"codigoVendedor"`
	ValorDescontoVenda   float64                   `json:"valorDescontoVenda"`
	ValorAcrescimoVenda  float64                   `json:"valorAcrescimoVenda"`
	Itens                []GetLancamentoPDVItemDTO `json:"itens"`
}

type GetLancamentoPDVItemDTO struct {
	CodigoBarras   string  `json:"codigoBarras"`
	Quantidade     float64 `json:"quantidade"`
	PrecoVenda     float64 `json:"precoVenda"`
	ValorDesconto  float64 `json:"valorDesconto"`
	ValorAcrescimo float64 `json:"valorAcrescimo"`
}

type UpdadeLancamentoPDVRequest struct {
	IDLoja     int   `json:"id_loja" binding:"required"`
	IDComanda  []int `json:"id_comanda" binding:"required"`
	Finalizado *bool `json:"finalizado" binding:"required"`
}
