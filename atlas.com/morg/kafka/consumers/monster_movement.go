package consumers

import (
	"atlas-morg/kafka/handler"
	"atlas-morg/monster"
	"github.com/sirupsen/logrus"
)

type monsterMovementEvent struct {
	UniqueId      int   `json:"uniqueId"`
	ObserverId    int   `json:"observerId"`
	SkillPossible bool  `json:"skillPossible"`
	Skill         int   `json:"skill"`
	SkillId       int   `json:"skillId"`
	SkillLevel    int   `json:"skillLevel"`
	Option        int   `json:"option"`
	StartX        int   `json:"startX"`
	StartY        int   `json:"startY"`
	EndX          int   `json:"endX"`
	EndY          int   `json:"endY"`
	Stance        int   `json:"stance"`
	RawMovement   []int `json:"rawMovement"`
}

func MonsterMovementEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &monsterMovementEvent{}
	}
}

func HandleMonsterMovementEvent() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*monsterMovementEvent); ok {
			monster.GetMonsterRegistry().MoveMonster(event.UniqueId, event.EndX, event.EndY, event.Stance)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
