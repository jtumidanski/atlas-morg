package monster

import (
	"atlas-morg/kafka/producers"
	_map "atlas-morg/map"
	attributes2 "atlas-morg/rest/attributes"
	requests2 "atlas-morg/rest/requests"
	"errors"
	"github.com/sirupsen/logrus"
)

func CreateMonster(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, input *attributes2.MonsterAttributes) (*Model, error) {
	return func(worldId byte, channelId byte, mapId uint32, input *attributes2.MonsterAttributes) (*Model, error) {
		ma, err := requests2.NewMapInformation(l).GetMonsterInformation(input.MonsterId)
		if err != nil {
			return nil, err
		}

		model := GetMonsterRegistry().CreateMonster(worldId, channelId, mapId, input.MonsterId, input.X, input.Y, input.Fh, 5, input.Team, ma.HP, ma.MP)

		cid, err := GetControllerCandidate(l)(worldId, channelId, mapId)
		if err == nil {
			StartControl(l)(model, cid)
		}

		producers.MonsterCreated(l)(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
		return model, nil
	}
}

func FindNextController(l logrus.FieldLogger) func(model *Model) {
	return func(model *Model) {
		cid, err := GetControllerCandidate(l)(model.WorldId(), model.ChannelId(), model.MapId())
		if err == nil {
			StartControl(l)(model, cid)
		}
	}
}

func GetControllerCandidate(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32) (uint32, error) {
		ids, err := _map.GetCharacterIdsInMap(l)(worldId, channelId, mapId)
		if err != nil {
			return 0, err
		}

		var controlCounts map[uint32]int
		controlCounts = make(map[uint32]int)

		for _, id := range ids {
			controlCounts[id] = 0
		}

		ms := GetMonsterRegistry().GetMonstersInMap(worldId, channelId, mapId)
		for _, m := range ms {
			if m.ControlCharacterId() != 0 {
				controlCounts[m.ControlCharacterId()] += 1
			}
		}

		var index = uint32(0)
		for key, val := range controlCounts {
			if index == 0 {
				index = key
			} else if val < controlCounts[index] {
				index = key
			}
		}

		if index == 0 {
			return 0, errors.New("should not get here")
		} else {
			return index, nil
		}
	}
}

func StartControl(l logrus.FieldLogger) func(monster *Model, controllerId uint32) {
	return func(monster *Model, controllerId uint32) {
		if monster.ControlCharacterId() != 0 {
			StopControl(l)(monster)
		}
		um := GetMonsterRegistry().ControlMonster(monster.UniqueId(), controllerId)
		if um != nil {
			producers.StartMonsterControl(l)(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
		}
	}
}

func StopControl(l logrus.FieldLogger) func(monster *Model) {
	return func(monster *Model) {
		um := GetMonsterRegistry().ClearControl(monster.UniqueId())
		if um != nil {
			producers.StopMonsterControl(l)(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
		}
	}
}

func DestroyAll(l logrus.FieldLogger) {
	ms := GetMonsterRegistry().GetMonsters()
	for _, x := range ms {
		Destroy(l)(x)
	}
}

func Destroy(l logrus.FieldLogger) func(model *Model) {
	return func(model *Model) {
		GetMonsterRegistry().RemoveMonster(model.UniqueId())
		producers.MonsterDestroyed(l)(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
	}
}
