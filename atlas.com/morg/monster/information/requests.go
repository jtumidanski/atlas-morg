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

func requestById(monsterId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(monsterResource, monsterId))
}
