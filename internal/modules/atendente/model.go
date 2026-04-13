package atendente

import "go.mongodb.org/mongo-driver/bson/primitive"

type Atendente struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IDLoja int                `bson:"idLoja" json:"idLoja"`
	Codigo string             `bson:"codigo" json:"codigo"`
	Nome   string             `bson:"nome" json:"nome"`
	Senha  string             `bson:"senha" json:"senha"`
	Ativo  bool               `bson:"ativo" json:"ativo"`
}
