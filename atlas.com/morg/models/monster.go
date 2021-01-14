package models

import "math"

type Monster struct {
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
	damageEntries      []DamageEntry
}

func NewMonster(worldId byte, channelId byte, mapId int, uniqueId int, monsterId int, x int, y int, fh int, stance int, team int, hp int) *Monster {
	return &Monster{
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
		damageEntries:      make([]DamageEntry, 0),
	}
}

func (m *Monster) UniqueId() int {
	return m.uniqueId
}

func (m *Monster) WorldId() byte {
	return m.worldId
}

func (m *Monster) ChannelId() byte {
	return m.channelId
}

func (m *Monster) MapId() int {
	return m.mapId
}

func (m *Monster) Hp() int {
	return m.hp
}

func (m *Monster) MonsterId() int {
	return m.monsterId
}

func (m *Monster) ControlCharacterId() int {
	return m.controlCharacterId
}

func (m *Monster) Fh() int {
	return m.fh
}

func (m *Monster) Team() int {
	return m.team
}

func (m *Monster) X() int {
	return m.x
}

func (m *Monster) Y() int {
	return m.y
}

func (m *Monster) Stance() int {
	return m.stance
}

func (m *Monster) DamageEntries() []DamageEntry {
	return m.damageEntries
}

func (m *Monster) Move(x int, y int, stance int) *Monster {
	return &Monster{
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

func (m *Monster) Control(characterId int) *Monster {
	return &Monster{
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

func (m *Monster) ClearControl() *Monster {
	return &Monster{
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

func (m *Monster) Damage(characterId int, damage int32) *Monster {
	actualDamage := int32(m.Hp()) - int32(math.Max(float64(m.Hp())-float64(damage), 0))

	return &Monster{
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
		damageEntries: append(m.DamageEntries(), DamageEntry{
			CharacterId: characterId,
			Damage:      actualDamage,
		}),
	}
}
