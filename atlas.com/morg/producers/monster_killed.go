package producers

import (
	"atlas-morg/events"
	"atlas-morg/models"
	"atlas-morg/requests"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type MonsterKilled struct {
	l   *log.Logger
	ctx context.Context
}

func NewMonsterKilled(l *log.Logger, ctx context.Context) *MonsterKilled {
	return &MonsterKilled{l, ctx}
}

func (mk *MonsterKilled) EmitKilled(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int, x int, y int, killerId int, damageSummary []models.DamageEntry) {
	t := requests.NewTopic(mk.l)
	td, err := t.GetTopic("TOPIC_MONSTER_KILLED_EVENT")
	if err != nil {
		mk.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        td.Attributes.Name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	var damageEntries []events.DamageEntry
	for _, x := range damageSummary {
		damageEntries = append(damageEntries, events.DamageEntry{
			Character: x.CharacterId,
			Damage:    x.Damage,
		})
	}

	e := &events.MonsterKilledEvent{
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
	r, err := json.Marshal(e)
	if err != nil {
		mk.l.Fatal("[ERROR] Unable to marshall event.")
	}

	mk.l.Printf("[INFO] Sending [MonsterKilledEvent] key %s", createKey(mapId))
	mk.l.Printf("[INFO] Sending [MonsterKilledEvent] value %s", r)

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   createKey(mapId),
		Value: r,
	})
	if err != nil {
		mk.l.Fatal("[ERROR] Unable to produce event.")
	}
}
