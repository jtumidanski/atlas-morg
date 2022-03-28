package monster

import (
	"atlas-morg/kafka/producers"
	_map "atlas-morg/map"
	"atlas-morg/monster/information"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func CreateMonster(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, input *Attributes) (*Model, error) {
	return func(worldId byte, channelId byte, mapId uint32, input *Attributes) (*Model, error) {
		ma, err := information.GetById(input.MonsterId)(l, span)
		if err != nil {
			return nil, err
		}

		attr := ma.Data().Attributes
		model := GetMonsterRegistry().CreateMonster(worldId, channelId, mapId, input.MonsterId, input.X, input.Y, input.Fh, 5, input.Team, attr.HP, attr.MP)

		cid, err := GetControllerCandidate(l, span)(worldId, channelId, mapId)
		if err == nil {
			StartControl(l, span)(model, cid)
		}

		producers.MonsterCreated(l, span)(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
		return model, nil
	}
}

func FindNextController(l logrus.FieldLogger, span opentracing.Span) func(model *Model) {
	return func(model *Model) {
		cid, err := GetControllerCandidate(l, span)(model.WorldId(), model.ChannelId(), model.MapId())
		if err == nil {
			StartControl(l, span)(model, cid)
		}
	}
}

func GetControllerCandidate(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) (uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32) (uint32, error) {
		ids, err := _map.GetCharacterIdsInMap(l, span)(worldId, channelId, mapId)
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

func StartControl(l logrus.FieldLogger, span opentracing.Span) func(monster *Model, controllerId uint32) {
	return func(monster *Model, controllerId uint32) {
		if monster.ControlCharacterId() != 0 {
			StopControl(l, span)(monster)
		}
		um := GetMonsterRegistry().ControlMonster(monster.UniqueId(), controllerId)
		if um != nil {
			producers.StartMonsterControl(l, span)(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
		}
	}
}

func StopControl(l logrus.FieldLogger, span opentracing.Span) func(monster *Model) {
	return func(monster *Model) {
		um := GetMonsterRegistry().ClearControl(monster.UniqueId())
		if um != nil {
			producers.StopMonsterControl(l, span)(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
		}
	}
}

func DestroyAll(l logrus.FieldLogger, span opentracing.Span) {
	ms := GetMonsterRegistry().GetMonsters()
	for _, x := range ms {
		Destroy(l, span)(x)
	}
}

func Destroy(l logrus.FieldLogger, span opentracing.Span) func(model *Model) {
	return func(model *Model) {
		GetMonsterRegistry().RemoveMonster(model.UniqueId())
		producers.MonsterDestroyed(l, span)(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
	}
}
