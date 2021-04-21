package monster

import (
	"atlas-morg/models"
	"math"
)

type Model struct {
	uniqueId           int
	worldId            byte
	channelId          byte
	mapId              int
	hp                 int
	monsterId          int
	controlCharacterId int
	x                  int
	y                  int
	fh                 int
	stance             int
	team               int
	damageEntries      []models.DamageEntry
}

func NewMonster(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int, x int, y int, fh int, stance int, team int, hp int) *Model {
	return &Model{
		uniqueId:           uniqueId,
		worldId:            worldId,
		channelId:          channelId,
		mapId:              mapId,
		hp:                 hp,
		monsterId:          monsterId,
		controlCharacterId: 0,
		x:                  x,
		y:                  y,
		fh:                 fh,
		stance:             stance,
		team:               team,
		damageEntries:      make([]models.DamageEntry, 0),
	}
}

func (m *Model) UniqueId() int {
	return m.uniqueId
}

func (m *Model) WorldId() byte {
	return m.worldId
}

func (m *Model) ChannelId() byte {
	return m.channelId
}

func (m *Model) MapId() int {
	return m.mapId
}

func (m *Model) Hp() int {
	return m.hp
}

func (m *Model) MonsterId() int {
	return m.monsterId
}

func (m *Model) ControlCharacterId() int {
	return m.controlCharacterId
}

func (m *Model) Fh() int {
	return m.fh
}

func (m *Model) Team() int {
	return m.team
}

func (m *Model) X() int {
	return m.x
}

func (m *Model) Y() int {
	return m.y
}

func (m *Model) Stance() int {
	return m.stance
}

func (m *Model) DamageEntries() []models.DamageEntry {
	return m.damageEntries
}

func (m *Model) DamageSummary() []models.DamageEntry {
	var damageSummary = make(map[int]int64)
	for _, x := range m.damageEntries {
		if _, ok := damageSummary[x.CharacterId]; ok {
			damageSummary[x.CharacterId] += x.Damage
		} else {
			damageSummary[x.CharacterId] = x.Damage
		}
	}

	var results []models.DamageEntry
	for id, dmg := range damageSummary {
		results = append(results, models.DamageEntry{
			CharacterId: id,
			Damage:      dmg,
		})
	}
	return results
}

func (m *Model) Move(x int, y int, stance int) *Model {
	return &Model{
		uniqueId:           m.UniqueId(),
		worldId:            m.WorldId(),
		channelId:          m.ChannelId(),
		mapId:              m.MapId(),
		hp:                 m.Hp(),
		monsterId:          m.MonsterId(),
		controlCharacterId: m.ControlCharacterId(),
		x:                  x,
		y:                  y,
		fh:                 m.Fh(),
		stance:             stance,
		team:               m.Team(),
		damageEntries:      m.DamageEntries(),
	}
}

func (m *Model) Control(characterId int) *Model {
	return &Model{
		uniqueId:           m.UniqueId(),
		worldId:            m.WorldId(),
		channelId:          m.ChannelId(),
		mapId:              m.MapId(),
		hp:                 m.Hp(),
		monsterId:          m.MonsterId(),
		controlCharacterId: characterId,
		x:                  m.X(),
		y:                  m.Y(),
		fh:                 m.Fh(),
		stance:             m.Stance(),
		team:               m.Team(),
		damageEntries:      m.DamageEntries(),
	}
}

func (m *Model) ClearControl() *Model {
	return &Model{
		uniqueId:           m.UniqueId(),
		worldId:            m.WorldId(),
		channelId:          m.ChannelId(),
		mapId:              m.MapId(),
		hp:                 m.Hp(),
		monsterId:          m.MonsterId(),
		controlCharacterId: 0,
		x:                  m.X(),
		y:                  m.Y(),
		fh:                 m.Fh(),
		stance:             m.Stance(),
		team:               m.Team(),
		damageEntries:      m.DamageEntries(),
	}
}

func (m *Model) Damage(characterId int, damage int64) *Model {
	actualDamage := int64(m.Hp()) - int64(math.Max(float64(m.Hp())-float64(damage), 0))

	return &Model{
		uniqueId:           m.UniqueId(),
		worldId:            m.WorldId(),
		channelId:          m.ChannelId(),
		mapId:              m.MapId(),
		hp:                 m.Hp() - int(actualDamage),
		monsterId:          m.MonsterId(),
		controlCharacterId: m.ControlCharacterId(),
		x:                  m.X(),
		y:                  m.Y(),
		fh:                 m.Fh(),
		stance:             m.Stance(),
		team:               m.Team(),
		damageEntries: append(m.DamageEntries(), models.DamageEntry{
			CharacterId: characterId,
			Damage:      actualDamage,
		}),
	}
}

func (m *Model) Alive() bool {
	return m.Hp() > 0
}

func (m *Model) DamageLeader() int {
	index := -1
	for i, x := range m.damageEntries {
		if index == -1 {
			index = i
		} else if m.damageEntries[index].Damage < x.Damage {
			index = i
		}
	}
	return m.damageEntries[index].CharacterId
}
