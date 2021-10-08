package consumers

import (
	"atlas-morg/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	MapCharacterEvent    = "map_character_event"
	MonsterDamageEvent   = "monster_damage_event"
	MonsterMovementEvent = "monster_movement_event"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, name string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, name, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_MAP_CHARACTER_EVENT", MapCharacterEvent, MapCharacterEventCreator(), HandleMapCharacterEvent())
	cec("TOPIC_MONSTER_DAMAGE", MonsterDamageEvent, MonsterDamageEventCreator(), HandleMonsterDamageEvent())
	cec("TOPIC_MONSTER_MOVEMENT", MonsterMovementEvent, MonsterMovementEventCreator(), HandleMonsterMovementEvent())
}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, name, topicToken, "Monster Registry Service", emptyEventCreator, processor)
}
