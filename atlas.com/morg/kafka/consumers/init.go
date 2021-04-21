package consumers

import (
	"context"
	"github.com/sirupsen/logrus"
)

func CreateEventConsumers(l *logrus.Logger) {
	cec := func(topicToken string, emptyEventCreator EmptyEventCreator, processor EventProcessor) {
		createEventConsumer(l, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_MAP_CHARACTER_EVENT", MapCharacterEventCreator(), HandleMapCharacterEvent())
	cec("TOPIC_MONSTER_DAMAGE", MonsterDamageEventCreator(), HandleMonsterDamageEvent())
	cec("TOPIC_MONSTER_MOVEMENT", MonsterMovementEventCreator(), HandleMonsterMovementEvent())
}

func createEventConsumer(l *logrus.Logger, topicToken string, emptyEventCreator EmptyEventCreator, processor EventProcessor) {
	h := func(logger logrus.FieldLogger, event interface{}) {
		processor(logger, event)
	}

	c := NewConsumer(l, context.Background(), h,
		SetGroupId("Monster Registry Service"),
		SetTopicToken(topicToken),
		SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}
