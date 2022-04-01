package information

import (
	"atlas-morg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetById(l logrus.FieldLogger, span opentracing.Span) func(monsterId uint32) (Model, error) {
	return func(monsterId uint32) (Model, error) {
		return requests.Provider[attributes, Model](l, span)(requestById(monsterId), makeModel)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	attr := body.Attributes
	return Model{
		hp: attr.HP,
		mp: attr.MP,
	}, nil
}
