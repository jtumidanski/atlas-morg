package consumers

import (
	"atlas-morg/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_MAP_CHARACTER_EVENT", MapCharacterEventCreator(), HandleMapCharacterEvent())
	cec("TOPIC_MONSTER_DAMAGE", MonsterDamageEventCreator(), HandleMonsterDamageEvent())
	cec("TOPIC_MONSTER_MOVEMENT", MonsterMovementEventCreator(), HandleMonsterMovementEvent())
}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, topicToken, "Monster Registry Service", emptyEventCreator, processor)
}
