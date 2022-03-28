package information

import (
	"atlas-morg/rest/requests"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	monstersResource                   = mapInformationService + "monsters"
	monsterResource                    = monstersResource + "/%d"
)

func GetById(monsterId uint32) requests.Request[MonsterDataAttributes] {
	return requests.MakeGetRequest[MonsterDataAttributes](fmt.Sprintf(monsterResource, monsterId))
}
