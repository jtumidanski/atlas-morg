package monster

import (
	"atlas-morg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameDamage   = "monster_damage_event"
	consumerNameMovement = "monster_movement_event"
	topicTokenDamage     = "TOPIC_MONSTER_DAMAGE"
	topicTokenMovement   = "TOPIC_MONSTER_MOVEMENT"
)

func DamageConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[damageEvent](consumerNameDamage, topicTokenDamage, groupId, handleDamage())
}

type damageEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	UniqueId    uint32 `json:"uniqueId"`
	CharacterId uint32 `json:"characterId"`
	Damage      int64  `json:"damage"`
}

func handleDamage() kafka.HandlerFunc[damageEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event damageEvent) {
		Damage(l, span)(event.UniqueId, event.CharacterId, event.Damage)
	}
}

func MovementConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[movementEvent](consumerNameMovement, topicTokenMovement, groupId, handleMovement())
}

type movementEvent struct {
	UniqueId      uint32 `json:"uniqueId"`
	ObserverId    int    `json:"observerId"`
	SkillPossible bool   `json:"skillPossible"`
	Skill         int    `json:"skill"`
	SkillId       int    `json:"skillId"`
	SkillLevel    int    `json:"skillLevel"`
	Option        int    `json:"option"`
	StartX        int    `json:"startX"`
	StartY        int    `json:"startY"`
	EndX          int    `json:"endX"`
	EndY          int    `json:"endY"`
	Stance        int    `json:"stance"`
	RawMovement   []int  `json:"rawMovement"`
}

func handleMovement() kafka.HandlerFunc[movementEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event movementEvent) {
		Move(l, span)(event.UniqueId, event.EndX, event.EndY, event.Stance)
	}
}
