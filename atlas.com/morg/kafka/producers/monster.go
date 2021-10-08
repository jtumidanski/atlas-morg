package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type monsterEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  int    `json:"uniqueId"`
	MonsterId int    `json:"monsterId"`
	ActorId   int    `json:"actorId"`
	Type      string `json:"type"`
}

func MonsterCreated(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int) {
	return func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int) {
		emitMonsterEvent(l, span)(worldId, channelId, mapId, uniqueId, monsterId, 0, "CREATED")
	}
}

func MonsterDestroyed(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int) {
	return func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int) {
		emitMonsterEvent(l, span)(worldId, channelId, mapId, uniqueId, monsterId, 0, "DESTROYED")
	}
}

func emitMonsterEvent(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int, actorId int, theType string) {
	return func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int, actorId int, theType string) {
		e := &monsterEvent{
			WorldId:   worldId,
			ChannelId: channelId,
			MapId:     mapId,
			UniqueId:  uniqueId,
			MonsterId: monsterId,
			ActorId:   actorId,
			Type:      theType,
		}
		ProduceEvent(l, span,"TOPIC_MONSTER_EVENT")(CreateKey(int(mapId)), e)
	}
}
