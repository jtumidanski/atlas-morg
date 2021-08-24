package consumers

import (
	"atlas-morg/kafka/handler"
	"atlas-morg/monster"
	"github.com/sirupsen/logrus"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func MapCharacterEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &mapCharacterEvent{}
	}
}

func HandleMapCharacterEvent() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*mapCharacterEvent); ok {
			if event.Type == "ENTER" {
				ms := monster.GetMonsterRegistry().GetMonstersInMap(event.WorldId, event.ChannelId, event.MapId)
				for _, m := range ms {
					if m.ControlCharacterId() == 0 {
						monster.FindNextController(l)(m)
					}
				}
			} else if event.Type == "EXIT" {
				ms := monster.GetMonsterRegistry().GetMonstersInMap(event.WorldId, event.ChannelId, event.MapId)
				for _, m := range ms {
					if m.ControlCharacterId() == event.CharacterId {
						monster.StopControl(l)(m)
					}
				}
				for _, m := range ms {
					if m.ControlCharacterId() == event.CharacterId {
						monster.FindNextController(l)(m)
					}
				}
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
