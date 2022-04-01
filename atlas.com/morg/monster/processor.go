package monster

import (
	_map "atlas-morg/map"
	"atlas-morg/monster/information"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func CreateMonster(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, input *Attributes) (*Model, error) {
	return func(worldId byte, channelId byte, mapId uint32, input *Attributes) (*Model, error) {
		ma, err := information.GetById(l, span)(input.MonsterId)
		if err != nil {
			return nil, err
		}
		model := GetMonsterRegistry().CreateMonster(worldId, channelId, mapId, input.MonsterId, input.X, input.Y, input.Fh, 5, input.Team, ma.HP(), ma.MP())

		cid, err := GetControllerCandidate(l, span)(worldId, channelId, mapId)
		if err == nil {
			StartControl(l, span)(model, cid)
		}

		emitCreated(l, span)(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
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
			emitStartControl(l, span)(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
		}
	}
}

func StopControl(l logrus.FieldLogger, span opentracing.Span) func(monster *Model) {
	return func(monster *Model) {
		um := GetMonsterRegistry().ClearControl(monster.UniqueId())
		if um != nil {
			emitStopControl(l, span)(um.WorldId(), um.ChannelId(), um.MapId(), um.ControlCharacterId(), um.UniqueId())
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
		emitDestroyed(l, span)(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
	}
}

func Move(_ logrus.FieldLogger, _ opentracing.Span) func(id uint32, x int, y int, stance int) {
	return func(id uint32, x int, y int, stance int) {
		GetMonsterRegistry().MoveMonster(id, x, y, stance)
	}
}

func Damage(l logrus.FieldLogger, span opentracing.Span) func(id uint32, characterId uint32, damage int64) {
	return func(id uint32, characterId uint32, damage int64) {
		m, err := GetMonsterRegistry().GetMonster(id)
		if err != nil {
			l.WithError(err).Errorf("Unable to get monster %d.", id)
			return
		}
		if !m.Alive() {
			l.Errorf("Character %d trying to apply damage to an already dead monster %d.", characterId, id)
			return
		}

		s, err := GetMonsterRegistry().ApplyDamage(characterId, damage, m.UniqueId())
		if err != nil {
			l.WithError(err).Errorf("Error applying damage to monster %d from character %d.", m.UniqueId(), characterId)
			return
		}

		if s.Killed {
			MonsterKilled(l, span)(s.Monster.WorldId(), s.Monster.ChannelId(), s.Monster.MapId(), s.Monster.UniqueId(), s.Monster.MonsterId(), s.Monster.X(), s.Monster.Y(), s.CharacterId, s.Monster.DamageSummary())
			GetMonsterRegistry().RemoveMonster(s.Monster.UniqueId())
			return
		}

		if characterId != s.Monster.ControlCharacterId() {
			dl := s.Monster.DamageLeader() == characterId
			if dl {
				StopControl(l, span)(&s.Monster)
				StartControl(l, span)(&s.Monster, characterId)
			}
		}

		// TODO broadcast HP bar update
	}
}
