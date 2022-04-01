package monster

import (
	"atlas-morg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type monsterEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
	MonsterId uint32 `json:"monsterId"`
	ActorId   int    `json:"actorId"`
	Type      string `json:"type"`
}

func emitCreated(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32) {
		emitEvent(l, span)(worldId, channelId, mapId, uniqueId, monsterId, 0, "CREATED")
	}
}

func emitDestroyed(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32) {
		emitEvent(l, span)(worldId, channelId, mapId, uniqueId, monsterId, 0, "DESTROYED")
	}
}

func emitEvent(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32, actorId int, theType string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MONSTER_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32, actorId int, theType string) {
		e := &monsterEvent{
			WorldId:   worldId,
			ChannelId: channelId,
			MapId:     mapId,
			UniqueId:  uniqueId,
			MonsterId: monsterId,
			ActorId:   actorId,
			Type:      theType,
		}
		producer(kafka.CreateKey(int(mapId)), e)
	}
}

type controlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	UniqueId    uint32 `json:"uniqueId"`
	Type        string `json:"type"`
}

func emitStartControl(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId uint32) {
		emitControl(l, span)(worldId, channelId, mapId, characterId, uniqueId, "START")
	}
}

func emitStopControl(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId uint32) {
		emitControl(l, span)(worldId, channelId, mapId, characterId, uniqueId, "STOP")
	}
}

func emitControl(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId uint32, theType string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CONTROL_MONSTER_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, uniqueId uint32, theType string) {
		e := &controlEvent{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
			UniqueId:    uniqueId,
			Type:        theType,
		}
		producer(kafka.CreateKey(int(mapId)), e)
	}
}

type killedEvent struct {
	WorldId       byte          `json:"worldId"`
	ChannelId     byte          `json:"channelId"`
	MapId         uint32        `json:"mapId"`
	UniqueId      uint32        `json:"uniqueId"`
	MonsterId     uint32        `json:"monsterId"`
	X             int           `json:"x"`
	Y             int           `json:"y"`
	KillerId      uint32        `json:"killerId"`
	DamageEntries []damageEntry `json:"damageEntries"`
}

type damageEntry struct {
	Character uint32 `json:"character"`
	Damage    int64  `json:"damage"`
}

func MonsterKilled(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32, x int, y int, killerId uint32, damageSummary []entry) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MONSTER_KILLED_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32, x int, y int, killerId uint32, damageSummary []entry) {
		var damageEntries []damageEntry
		for _, x := range damageSummary {
			damageEntries = append(damageEntries, damageEntry{
				Character: x.CharacterId,
				Damage:    x.Damage,
			})
		}
		e := &killedEvent{
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
		producer(kafka.CreateKey(int(mapId)), e)
	}
}
