package information

import (
	"atlas-morg/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	monstersResource                   = mapInformationService + "monsters"
	monsterResource                    = monstersResource + "/%d"
)

func GetById(l logrus.FieldLogger, span opentracing.Span) func(monsterId uint32) (*MonsterDataData, error) {
	return func(monsterId uint32) (*MonsterDataData, error) {
		td := &MonsterDataDataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(monsterResource, monsterId), td)
		if err != nil {
			return nil, err
		}
		return &td.Data, nil
	}
}
