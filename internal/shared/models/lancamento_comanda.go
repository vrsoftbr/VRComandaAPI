package models

import "time"

type LancamentoComanda struct {
	ID          uint                    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDLoja      int                     `gorm:"column:id_loja;not null;index" json:"id_loja"`
	IDComanda   int                     `gorm:"column:id_comanda;not null;index" json:"id_comanda"`
	IDMesa      *int                    `gorm:"column:id_mesa" json:"id_mesa,omitempty"`
	IDAtendente int                     `gorm:"column:id_atendente;not null;index" json:"id_atendente"`
	DataHora    time.Time               `gorm:"column:data_hora;not null;index" json:"dataHora"`
	Observacao  string                  `gorm:"column:observacao;size:255" json:"observacao,omitempty"`
	Finalizado  bool                    `gorm:"column:finalizado;not null;default:false;index" json:"finalizado"`
	Itens       []LancamentoComandaItem `gorm:"foreignKey:IDLancamentoComanda;constraint:OnDelete:CASCADE" json:"itens,omitempty"`
}

func (LancamentoComanda) TableName() string {
	return "lancamentocomanda"
}

type LancamentoComandaItem struct {
	ID                   uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDLancamentoComanda  uint       `gorm:"column:id_lancamentocomanda;not null;index" json:"id_lancamentocomanda"`
	Sequencia            int        `gorm:"column:sequencia;not null" json:"sequencia"`
	IDProduto            int        `gorm:"column:id_produto;not null" json:"id_produto"`
	CodigoBarras         string     `gorm:"column:codigobarras;size:64" json:"codigobarras,omitempty"`
	Quantidade           float64    `gorm:"column:quantidade;not null" json:"quantidade"`
	PrecoVenda           float64    `gorm:"column:precovenda;not null" json:"precovenda"`
	Cancelado            bool       `gorm:"column:cancelado;not null;default:false" json:"cancelado"`
	DataHoraCancelamento *time.Time `gorm:"column:data_hora_cancelamento" json:"dataHoraCancelamento,omitempty"`
	IDAtendente          int        `gorm:"column:id_atendente;not null;index" json:"id_atendente"`
	IDSituacao           int        `gorm:"column:id_situacao;not null;default:0" json:"id_situacao"`
}

func (LancamentoComandaItem) TableName() string {
	return "lancamentocomandaitem"
}
