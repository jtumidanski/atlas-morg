package consumers

import (
	"atlas-morg/events"
	"atlas-morg/registries"
	"atlas-morg/requests"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type MonsterMovement struct {
	l   *log.Logger
	ctx context.Context
}

func NewMonsterMovement(l *log.Logger, ctx context.Context) *MonsterMovement {
	return &MonsterMovement{l, ctx}
}

func (mc *MonsterMovement) Init() {
	t := requests.NewTopic(mc.l)
	td, err := t.GetTopic("TOPIC_MONSTER_MOVEMENT")
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

		var event events.MonsterMovementEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			mc.l.Println("Could not unmarshal event into event class ", msg.Value)
			mc.l.Println(err.Error())
		} else {
			mc.processEvent(event)
		}
	}
}

func (mc *MonsterMovement) processEvent(event events.MonsterMovementEvent) {
	registries.GetMonsterRegistry().MoveMonster(event.UniqueId, event.EndX, event.EndY, event.Stance)
}
