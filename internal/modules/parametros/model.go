package parametros

import "go.mongodb.org/mongo-driver/bson/primitive"

type Parametro struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	IDParametro int                `bson:"idParametro" json:"idParametro"`
	Descricao   string             `bson:"descricao" json:"descricao"`
	IDLoja      int                `bson:"idLoja" json:"idLoja"`
	Valor       string             `bson:"valor" json:"valor"`
	Classe      string             `bson:"_class,omitempty" json:"_class,omitempty"`
}
