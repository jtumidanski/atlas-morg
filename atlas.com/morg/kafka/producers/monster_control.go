package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type monsterControlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	UniqueId    int    `json:"uniqueId"`
	Type        string `json:"type"`
}

func StartMonsterControl(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId int) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId int) {
		emitMonsterControlEvent(l, span)(worldId, channelId, mapId, characterId, uniqueId, "START")
	}
}

func StopMonsterControl(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId int) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId int) {
		emitMonsterControlEvent(l, span)(worldId, channelId, mapId, characterId, uniqueId, "STOP")
	}
}

func emitMonsterControlEvent(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId int, theType string) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId int, theType string) {
		e := &monsterControlEvent{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
			UniqueId:    uniqueId,
			Type:        theType,
		}
		ProduceEvent(l, span, "TOPIC_CONTROL_MONSTER_EVENT")(CreateKey(int(mapId)), e)
	}
}
