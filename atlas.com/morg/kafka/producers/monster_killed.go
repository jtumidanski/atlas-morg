package producers

import (
	"atlas-morg/models"
	"github.com/sirupsen/logrus"
)

type monsterKilledEvent struct {
	WorldId       byte          `json:"worldId"`
	ChannelId     byte          `json:"channelId"`
	MapId         uint32        `json:"mapId"`
	UniqueId      int           `json:"uniqueId"`
	MonsterId     int           `json:"monsterId"`
	X             int           `json:"x"`
	Y             int           `json:"y"`
	KillerId      uint32        `json:"killerId"`
	DamageEntries []damageEntry `json:"damageEntries"`
}

type damageEntry struct {
	Character uint32 `json:"character"`
	Damage    int64  `json:"damage"`
}

func MonsterKilled(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int, x int, y int, killerId uint32, damageSummary []models.DamageEntry) {
	return func(worldId byte, channelId byte, mapId uint32, uniqueId int, monsterId int, x int, y int, killerId uint32, damageSummary []models.DamageEntry) {
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
		ProduceEvent(l, "TOPIC_MONSTER_KILLED_EVENT")(CreateKey(int(mapId)), e)
	}
}
