package information

import (
	"atlas-morg/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	monstersResource                   = mapInformationService + "monsters"
	monsterResource                    = monstersResource + "/%d"
)

func GetById(l logrus.FieldLogger) func(monsterId int) (*MonsterDataData, error) {
	return func(monsterId int) (*MonsterDataData, error) {
		td := &MonsterDataDataContainer{}
		err := requests.Get(l)(fmt.Sprintf(monsterResource, monsterId), td)
		if err != nil {
			return nil, err
		}
		return &td.Data, nil
	}
}
