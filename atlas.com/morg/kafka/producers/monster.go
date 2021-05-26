package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type monsterEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     int    `json:"mapId"`
	UniqueId  int    `json:"uniqueId"`
	MonsterId int    `json:"monsterId"`
	ActorId   int    `json:"actorId"`
	Type      string `json:"type"`
}

type Monster struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func NewMonster(l logrus.FieldLogger, ctx context.Context) *Monster {
	return &Monster{l, ctx}
}

func (m *Monster) EmitCreated(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int) {
	m.emit(worldId, channelId, mapId, uniqueId, monsterId, 0, "CREATED")
}

func (m *Monster) EmitDestroyed(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int) {
	m.emit(worldId, channelId, mapId, uniqueId, monsterId, 0, "DESTROYED")
}

func (m *Monster) emit(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int, actorId int, theType string) {
	e := &monsterEvent{
		WorldId:   worldId,
		ChannelId: channelId,
		MapId:     mapId,
		UniqueId:  uniqueId,
		MonsterId: monsterId,
		ActorId:   actorId,
		Type:      theType,
	}
	ProduceEvent(m.l, "TOPIC_MONSTER_EVENT")(CreateKey(mapId), e)
}
