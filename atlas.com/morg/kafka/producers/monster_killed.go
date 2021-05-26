package producers

import (
	"atlas-morg/models"
	"context"
	"github.com/sirupsen/logrus"
)

type monsterKilledEvent struct {
	WorldId       byte          `json:"worldId"`
	ChannelId     byte          `json:"channelId"`
	MapId         int           `json:"mapId"`
	UniqueId      int           `json:"uniqueId"`
	MonsterId     int           `json:"monsterId"`
	X             int           `json:"x"`
	Y             int           `json:"y"`
	KillerId      int           `json:"killerId"`
	DamageEntries []damageEntry `json:"damageEntries"`
}

type damageEntry struct {
	Character int   `json:"character"`
	Damage    int64 `json:"damage"`
}

type MonsterKilled struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func NewMonsterKilled(l logrus.FieldLogger, ctx context.Context) *MonsterKilled {
	return &MonsterKilled{l, ctx}
}

func (mk *MonsterKilled) EmitKilled(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int, x int, y int, killerId int, damageSummary []models.DamageEntry) {
	var damageEntries []damageEntry
	for _, x := range damageSummary {
		damageEntries = append(damageEntries, damageEntry{
			Character: x.CharacterId,
			Damage:    x.Damage,
		})
	}
	e := &monsterKilledEvent{
		WorldId:       worldId,
		ChannelId:     channelId,
		MapId:         mapId,
		UniqueId:      uniqueId,
		MonsterId:     monsterId,
		X:             x,
		Y:             y,
		KillerId:      killerId,
		DamageEntries: damageEntries,
	}
	ProduceEvent(mk.l, "TOPIC_MONSTER_KILLED_EVENT")(CreateKey(mapId), e)
}
