package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type monsterControlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId int    `json:"characterId"`
	UniqueId    int    `json:"uniqueId"`
	Type        string `json:"type"`
}

type MonsterControl struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func NewMonsterControl(l logrus.FieldLogger, ctx context.Context) *MonsterControl {
	return &MonsterControl{l, ctx}
}

func (m *MonsterControl) EmitControl(worldId byte, channelId byte, mapId int, characterId int, uniqueId int) {
	m.emit(worldId, channelId, mapId, characterId, uniqueId, "START")
}

func (m *MonsterControl) EmitStop(worldId byte, channelId byte, mapId int, characterId int, uniqueId int) {
	m.emit(worldId, channelId, mapId, characterId, uniqueId, "STOP")
}

func (m *MonsterControl) emit(worldId byte, channelId byte, mapId int, characterId int, uniqueId int, theType string) {
	e := &monsterControlEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		CharacterId: characterId,
		UniqueId:    uniqueId,
		Type:        theType,
	}
	ProduceEvent(m.l, "TOPIC_CONTROL_MONSTER_EVENT")(CreateKey(mapId), e)
}
