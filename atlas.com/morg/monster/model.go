package monster

import (
	"atlas-morg/models"
	"math"
)

type Model struct {
	uniqueId           uint32
	worldId            byte
	channelId          byte
	mapId              uint32
	maxHp              uint32
	hp                 uint32
	maxMp              uint32
	mp                 uint32
	monsterId          uint32
	controlCharacterId uint32
	x                  int
	y                  int
	fh                 int
	stance             int
	team               int
	damageEntries      []models.DamageEntry
}

func NewMonster(worldId byte, channelId byte, mapId uint32, uniqueId uint32, monsterId uint32, x int, y int, fh int, stance int, team int, hp uint32, mp uint32) *Model {
	return &Model{
		uniqueId:           uniqueId,
		worldId:            worldId,
		channelId:          channelId,
		mapId:              mapId,
		maxHp:              hp,
		hp:                 hp,
		maxMp:              mp,
		mp:                 mp,
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

func (m *Model) UniqueId() uint32 {
	return m.uniqueId
}

func (m *Model) WorldId() byte {
	return m.worldId
}

func (m *Model) ChannelId() byte {
	return m.channelId
}

func (m *Model) MapId() uint32 {
	return m.mapId
}

func (m *Model) Hp() uint32 {
	return m.hp
}

func (m *Model) MonsterId() uint32 {
	return m.monsterId
}

func (m *Model) ControlCharacterId() uint32 {
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
	var damageSummary = make(map[uint32]int64)
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
		maxHp:              m.MaxHp(),
		hp:                 m.Hp(),
		maxMp:              m.MaxMp(),
		mp:                 m.Mp(),
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

func (m *Model) Control(characterId uint32) *Model {
	return &Model{
		uniqueId:           m.UniqueId(),
		worldId:            m.WorldId(),
		channelId:          m.ChannelId(),
		mapId:              m.MapId(),
		maxHp:              m.MaxHp(),
		hp:                 m.Hp(),
		maxMp:              m.MaxMp(),
		mp:                 m.Mp(),
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
		maxHp:              m.MaxHp(),
		hp:                 m.Hp(),
		maxMp:              m.MaxMp(),
		mp:                 m.Mp(),
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

func (m *Model) Damage(characterId uint32, damage int64) *Model {
	actualDamage := int64(m.Hp()) - int64(math.Max(float64(m.Hp())-float64(damage), 0))

	return &Model{
		uniqueId:           m.UniqueId(),
		worldId:            m.WorldId(),
		channelId:          m.ChannelId(),
		mapId:              m.MapId(),
		maxHp:              m.MaxHp(),
		hp:                 m.Hp() - uint32(actualDamage),
		maxMp:              m.MaxMp(),
		mp:                 m.Mp(),
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

func (m *Model) DamageLeader() uint32 {
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

func (m *Model) MaxHp() uint32 {
	return m.maxHp
}

func (m *Model) MaxMp() uint32 {
	return m.maxMp
}

func (m *Model) Mp() uint32 {
	return m.mp
}
