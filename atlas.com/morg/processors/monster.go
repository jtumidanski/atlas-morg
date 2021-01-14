package processor

import (
	"atlas-morg/attributes"
	"atlas-morg/models"
	"atlas-morg/registries"
	"errors"
	"log"
)

type Monster struct {
	l *log.Logger
}

func NewMonster(l *log.Logger) *Monster {
	return &Monster{l}
}

func (m *Monster) CreateMonster(worldId byte, channelId byte, mapId int, input *attributes.MonsterAttributes) (*models.Monster, error) {
	ma, err := NewMapInformation(m.l).GetMonsterInformation(input.MonsterId)
	if err != nil {
		return nil, err
	}

	model := registries.GetMonsterRegistry().CreateMonster(worldId, channelId, mapId, input.MonsterId, input.X, input.Y, input.Fh, 5, input.Team, ma.Hp)

	cid, err := m.GetControllerCandidate(worldId, channelId, mapId)
	if err == nil {
		m.StartControl(model, cid)
	}

	// emit monster created event
	return model, nil
}

func (m *Monster) GetControllerCandidate(worldId byte, channelId byte, mapId int) (int, error) {
	ids, err := NewMap(m.l).GetCharacterIdsInMap(worldId, channelId, mapId)
	if err != nil {
		return 0, err
	}

	var controlCounts map[int]int
	controlCounts = make(map[int]int)

	for _, id := range ids {
		controlCounts[id] = 0
	}

	ms := registries.GetMonsterRegistry().GetMonstersInMap(worldId, channelId, mapId)
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

func (m *Monster) StartControl(monster *models.Monster, controllerId int) {
	if monster.ControlCharacterId() != 0 {
		m.StopControl(monster)
	}
	registries.GetMonsterRegistry().ControlMonster(monster.UniqueId(), controllerId)
	// emit start control event
}

func (m *Monster) StopControl(monster *models.Monster) {
	registries.GetMonsterRegistry().ClearControl(monster.UniqueId())
	// emit clear control event
}
