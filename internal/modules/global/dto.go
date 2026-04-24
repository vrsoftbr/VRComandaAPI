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
	IDLoja     int   `form:"id_loja"`
	Finalizado *bool `form:"finalizado"`
}

type ListLancamentosDetalhesRequest struct {
	GlobalFilterRequest
}
