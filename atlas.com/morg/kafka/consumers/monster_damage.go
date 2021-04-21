package consumers

import (
	"atlas-morg/kafka/producers"
	"atlas-morg/monster"
	"context"
	"github.com/sirupsen/logrus"
)

type monsterDamageEvent struct {
	WorldId     byte  `json:"worldId"`
	ChannelId   byte  `json:"channelId"`
	MapId       int   `json:"mapId"`
	UniqueId    int   `json:"uniqueId"`
	CharacterId int   `json:"characterId"`
	Damage      int64 `json:"damage"`
}

func MonsterDamageEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &monsterDamageEvent{}
	}
}

func HandleMonsterDamageEvent() EventProcessor {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*monsterDamageEvent); ok {
			m, err := monster.GetMonsterRegistry().GetMonster(event.UniqueId)
			if err != nil {
				l.WithError(err).Errorf("Unable to get monster %d.", event.UniqueId)
				return
			}
			if !m.Alive() {
				l.Errorf("Character %d trying to apply damage to an already dead monster %d.", event.CharacterId, event.UniqueId)
				return
			}

			s, err := monster.GetMonsterRegistry().ApplyDamage(event.CharacterId, event.Damage, m.UniqueId())
			if err != nil {
				l.WithError(err).Errorf("Error applying damage to monster %d from character %d.", m.UniqueId(), event.CharacterId)
				return
			}

			if s.Killed {
				producers.NewMonsterKilled(l, context.Background()).EmitKilled(s.Monster.WorldId(), s.Monster.ChannelId(), s.Monster.MapId(), s.Monster.UniqueId(), s.Monster.MonsterId(), s.Monster.X(), s.Monster.Y(), s.CharacterId, s.Monster.DamageSummary())
				monster.GetMonsterRegistry().RemoveMonster(s.Monster.UniqueId())
				return
			}

			if event.CharacterId != s.Monster.ControlCharacterId() {
				dl := s.Monster.DamageLeader() == event.CharacterId
				if dl {
					monster.Processor(l).StopControl(&s.Monster)
					monster.Processor(l).StartControl(&s.Monster, event.CharacterId)
				}
			}

			// TODO broadcast HP bar update
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
