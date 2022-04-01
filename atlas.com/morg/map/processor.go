package _map

import (
	"atlas-morg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetCharacterIdsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
		return requests.SliceProvider[characterAttributes, uint32](l, span)(requestCharactersInMap(worldId, channelId, mapId), getCharacterId)()
	}
}

func getCharacterId(body requests.DataBody[characterAttributes]) (uint32, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}
