package mesa

import "go.mongodb.org/mongo-driver/bson/primitive"

type Mesa struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IDLoja    int                `bson:"idLoja" json:"idLoja"`
	Mesa      int                `bson:"mesa" json:"mesa"`
	Descricao string             `bson:"descricao" json:"descricao"`
	Ativo     bool               `bson:"ativo" json:"ativo"`
}
