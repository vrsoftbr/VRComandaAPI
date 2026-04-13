package comanda

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comanda struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IDLoja              int                `bson:"idLoja" json:"idLoja"`
	Comanda             int                `bson:"comanda" json:"comanda"`
	NumeroIdentificacao string             `bson:"numeroIdentificacao" json:"numeroIdentificacao"`
	Observacao          string             `bson:"observacao" json:"observacao"`
	Ativo               bool               `bson:"ativo" json:"ativo"`
}
