package producers

import (
	"atlas-morg/events"
	"atlas-morg/requests"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type Monster struct {
	l   *log.Logger
	ctx context.Context
}

func NewMonster(l *log.Logger, ctx context.Context) *Monster {
	return &Monster{l, ctx}
}

func (m *Monster) EmitCreated(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int) {
	m.emit(worldId, channelId, mapId, uniqueId, monsterId, 0, "CREATED")
}

func (m *Monster) EmitDestroyed(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int) {
	m.emit(worldId, channelId, mapId, uniqueId, monsterId, 0, "DESTROYED")
}

func (m *Monster) emit(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int, actorId int, theType string) {
	t := requests.NewTopic(m.l)
	td, err := t.GetTopic("TOPIC_MONSTER_EVENT")
	if err != nil {
		m.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        td.Attributes.Name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	e := &events.MonsterEvent{
		WorldId:   worldId,
		ChannelId: channelId,
		MapId:     mapId,
		UniqueId:  uniqueId,
		MonsterId: monsterId,
		ActorId:   actorId,
		Type:      theType,
	}
	r, err := json.Marshal(e)
	if err != nil {
		m.l.Fatal("[ERROR] Unable to marshall event.")
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   createKey(mapId),
		Value: r,
	})
	if err != nil {
		m.l.Print(err.Error())
		m.l.Fatal("[ERROR] Unable to produce monster event.")
	}
}
