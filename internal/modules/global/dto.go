package global

import (
	"vrcomandaapi/internal/modules/comanda"
	"vrcomandaapi/internal/modules/mesa"
	"vrcomandaapi/internal/shared/models"
)

type LancamentoDetalhesDTO struct {
	IDLancamento uint                           `json:"id_lancamento"`
	IDLoja       int                            `json:"id_loja"`
	IDComanda    int                            `json:"id_comanda"`
	IDMesa       *int                           `json:"id_mesa,omitempty"`
	DataHora     string                         `json:"dataHora"`
	Observacao   string                         `json:"observacao,omitempty"`
	IDAtendente  int                            `json:"id_atendente"`
	Finalizado   bool                           `json:"finalizado"`
	Comanda      *comanda.ComandaResponse       `json:"comanda,omitempty"`
	Mesa         *mesa.MesaResponse             `json:"mesa,omitempty"`
	Itens        []models.LancamentoComandaItem `json:"itens,omitempty"`
}

type GlobalFilterRequest struct {
	IDLoja     int   `form:"idLoja"`
	Finalizado *bool `form:"finalizado"`
}

type ListLancamentosDetalhesRequest struct {
	GlobalFilterRequest
}

type SituacaoComanda int

const (
	SituacaoComandaLiberada  SituacaoComanda = 1
	SituacaoComandaBloqueada SituacaoComanda = 2
)

type ConsultarSituacaoComandaRequest struct {
	IDLoja                     int    `form:"idLoja"`
	NumeroIdentificacaoComanda string `form:"numeroIdentificacaoComanda"`
}

type ConsultarSituacaoComandaResponse struct {
	IDLoja                     int             `json:"idLoja"`
	Comanda                    int             `json:"comanda"`
	NumeroIdentificacaoComanda string          `json:"numeroIdentificacaoComanda"`
	Situacao                   SituacaoComanda `json:"situacao"`
}
