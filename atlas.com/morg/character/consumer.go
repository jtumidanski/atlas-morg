package character

import (
	"atlas-morg/kafka"
	"atlas-morg/monster"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameMapCharacter = "map_character_event"
	topicTokenMapCharacter   = "TOPIC_MAP_CHARACTER_EVENT"
)

func MapConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[mapEvent](consumerNameMapCharacter, topicTokenMapCharacter, groupId, handleMapEvent())
}

type mapEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func handleMapEvent() kafka.HandlerFunc[mapEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event mapEvent) {
		if event.Type == "ENTER" {
			ms := monster.GetMonsterRegistry().GetMonstersInMap(event.WorldId, event.ChannelId, event.MapId)
			for _, m := range ms {
				if m.ControlCharacterId() == 0 {
					monster.FindNextController(l, span)(m)
				}
			}
		} else if event.Type == "EXIT" {
			ms := monster.GetMonsterRegistry().GetMonstersInMap(event.WorldId, event.ChannelId, event.MapId)
			for _, m := range ms {
				if m.ControlCharacterId() == event.CharacterId {
					monster.StopControl(l, span)(m)
				}
			}
			for _, m := range ms {
				if m.ControlCharacterId() == event.CharacterId {
					monster.FindNextController(l, span)(m)
				}
			}
		}
	}
}
