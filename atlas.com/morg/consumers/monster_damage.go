package consumers

import (
	"atlas-morg/events"
	"atlas-morg/models"
	"atlas-morg/processors"
	"atlas-morg/producers"
	"atlas-morg/registries"
	"atlas-morg/requests"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type MonsterDamage struct {
	l   *log.Logger
	ctx context.Context
}

func NewMonsterDamage(l *log.Logger, ctx context.Context) *MonsterDamage {
	return &MonsterDamage{l, ctx}
}

func (mc *MonsterDamage) Init() {
	t := requests.NewTopic(mc.l)
	td, err := t.GetTopic("TOPIC_MONSTER_DAMAGE")
	if err != nil {
		mc.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: "Monster Registry",
		MaxWait: 50 * time.Millisecond,
	})
	for {
		msg, err := r.ReadMessage(mc.ctx)
		if err != nil {
			panic("Could not successfully read message " + err.Error())
		}

		var event events.MonsterDamageEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			mc.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			mc.processEvent(event)
		}
	}
}

func (mc *MonsterDamage) processEvent(event events.MonsterDamageEvent) {
	m, err := registries.GetMonsterRegistry().GetMonster(event.UniqueId)
	if err == nil && m.Alive() {
		mc.applyDamage(event.CharacterId, event.Damage, m)
	}
}

func (mc *MonsterDamage) applyDamage(characterId int, damage int64, m *models.Monster) {
	s, err := registries.GetMonsterRegistry().ApplyDamage(characterId, damage, m.UniqueId())
	if err == nil {
		if s.Killed {
			mc.killMonster(s)
			return
		}

		if characterId != s.Monster.ControlCharacterId() {
			dl := s.Monster.DamageLeader() == characterId
			if dl {
				processors.NewMonster(mc.l).StopControl(&s.Monster)
				processors.NewMonster(mc.l).StartControl(&s.Monster, characterId)
			}
		}

		// TODO broadcast HP bar update
	}
}

func (mc *MonsterDamage) killMonster(s *models.DamageSummary) {
	producers.NewMonsterKilled(mc.l, mc.ctx).EmitKilled(s.Monster.WorldId(), s.Monster.ChannelId(), s.Monster.MapId(), s.Monster.UniqueId(), s.Monster.MonsterId(), s.Monster.X(), s.Monster.Y(), s.CharacterId, s.Monster.DamageSummary())
	registries.GetMonsterRegistry().RemoveMonster(s.Monster.UniqueId())
}
