package monster

import (
	producers2 "atlas-morg/kafka/producers"
	attributes2 "atlas-morg/rest/attributes"
	requests2 "atlas-morg/rest/requests"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

type processor struct {
	l logrus.FieldLogger
}

var Processor = func(l logrus.FieldLogger) *processor {
	return &processor{l}
}

func (m *processor) CreateMonster(worldId byte, channelId byte, mapId int, input *attributes2.MonsterAttributes) (*Model, error) {
	ma, err := requests2.NewMapInformation(m.l).GetMonsterInformation(input.MonsterId)
	if err != nil {
		return nil, err
	}

	model := GetMonsterRegistry().CreateMonster(worldId, channelId, mapId, input.MonsterId, input.X, input.Y, input.Fh, 5, input.Team, ma.Hp, ma.Mp)

	cid, err := m.GetControllerCandidate(worldId, channelId, mapId)
	if err == nil {
		m.StartControl(model, cid)
	}

	producers2.NewMonster(m.l, context.Background()).EmitCreated(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
	return model, nil
}

func (m *processor) FindNextController(model *Model) {
	cid, err := m.GetControllerCandidate(model.WorldId(), model.ChannelId(), model.MapId())
	if err == nil {
		m.StartControl(model, cid)
	}
}

func (m *processor) GetControllerCandidate(worldId byte, channelId byte, mapId int) (int, error) {
	ids, err := requests2.NewMap(m.l).GetCharacterIdsInMap(worldId, channelId, mapId)
	if err != nil {
		return 0, err
	}

	var controlCounts map[int]int
	controlCounts = make(map[int]int)

	for _, id := range ids {
		controlCounts[id] = 0
	}

	ms := GetMonsterRegistry().GetMonstersInMap(worldId, channelId, mapId)
	for _, m := range ms {
		if m.ControlCharacterId() != 0 {
			controlCounts[m.ControlCharacterId()] += 1
		}
	}

	var index = -1
	for key, val := range controlCounts {
		if index == -1 {
			index = key
		} else if val < controlCounts[index] {
			index = key
		}
	}

	if index == -1 {
		return 0, errors.New("should not get here")
	} else {
		return index, nil
	}
}

func (m *processor) StartControl(monster *Model, controllerId int) {
	if monster.ControlCharacterId() != 0 {
		m.StopControl(monster)
	}
	um := GetMonsterRegistry().ControlMonster(monster.UniqueId(), controllerId)
	if um != nil {
		producers2.NewMonsterControl(m.l, context.Background()).EmitControl(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
	}
}

func (m *processor) StopControl(monster *Model) {
	um := GetMonsterRegistry().ClearControl(monster.UniqueId())
	if um != nil {
		producers2.NewMonsterControl(m.l, context.Background()).EmitStop(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
	}
}

func (m *processor) DestroyAll() {
	ms := GetMonsterRegistry().GetMonsters()
	for _, x := range ms {
		m.Destroy(&x)
	}
}

func (m *processor) Destroy(model *Model) {
	GetMonsterRegistry().RemoveMonster(model.UniqueId())
	producers2.NewMonster(m.l, context.Background()).EmitDestroyed(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
}
