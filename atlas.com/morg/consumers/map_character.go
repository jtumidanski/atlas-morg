package consumers

import (
	"atlas-morg/events"
	"atlas-morg/processors"
	"atlas-morg/registries"
	"atlas-morg/requests"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type MapCharacter struct {
	l   *log.Logger
	ctx context.Context
}

func NewMapCharacter(l *log.Logger, ctx context.Context) *MapCharacter {
	return &MapCharacter{l, ctx}
}

func (mc *MapCharacter) Init() {
	t := requests.NewTopic(mc.l)
	td, err := t.GetTopic("TOPIC_MAP_CHARACTER_EVENT")
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

		var event events.MapCharacterEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			mc.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			mc.processEvent(event)
		}
	}
}

func (mc *MapCharacter) processEvent(event events.MapCharacterEvent) {
	if event.Type == "ENTER" {
		mc.gainControl(event)
	} else if event.Type == "EXIT" {
		mc.removeControl(event)
	}
}

func (mc *MapCharacter) gainControl(event events.MapCharacterEvent) {
	ms := registries.GetMonsterRegistry().GetMonstersInMap(event.WorldId, event.ChannelId, event.MapId)
	for _, m := range ms {
		if m.ControlCharacterId() == 0 {
			processors.NewMonster(mc.l).FindNextController(&m)
		}
	}
}

func (mc *MapCharacter) removeControl(event events.MapCharacterEvent) {
	ms := registries.GetMonsterRegistry().GetMonstersInMap(event.WorldId, event.ChannelId, event.MapId)
	for _, m := range ms {
		if m.ControlCharacterId() == event.CharacterId {
			processors.NewMonster(mc.l).StopControl(&m)
		}
	}
	for _, m := range ms {
		if m.ControlCharacterId() == event.CharacterId {
			processors.NewMonster(mc.l).FindNextController(&m)
		}
	}
}
