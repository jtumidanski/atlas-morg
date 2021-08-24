package monster

import (
	"atlas-morg/models"
	"errors"
	"sync"
)

type Registry struct {
	mutex                 sync.Mutex
	monsterRegisterRWLock sync.RWMutex
	monsterRegister       map[int]*Model
	mapMonsters           map[models.MapKey][]int
	mapLocks              map[models.MapKey]*sync.Mutex
}

var monsterRegistry *Registry
var once sync.Once

var uniqueId = 1000000001

func GetMonsterRegistry() *Registry {
	once.Do(func() {
		monsterRegistry = &Registry{}

		monsterRegistry.monsterRegister = make(map[int]*Model)
		monsterRegistry.mapMonsters = make(map[models.MapKey][]int)

		monsterRegistry.mapLocks = make(map[models.MapKey]*sync.Mutex)
	})
	return monsterRegistry
}

func (r *Registry) getMapLock(key models.MapKey) *sync.Mutex {
	if val, ok := r.mapLocks[key]; ok {
		return val
	} else {
		var cm = &sync.Mutex{}
		r.mutex.Lock()
		r.mapLocks[key] = cm
		r.mutex.Unlock()
		return cm
	}
}

func existingIds(monsters map[int]*Model) []int {
	var ids []int
	for _, x := range monsters {
		ids = append(ids, x.UniqueId())
	}
	return ids
}

func contains(ids []int, id int) bool {
	for _, element := range ids {
		if element == id {
			return true
		}
	}
	return false
}

func (r *Registry) CreateMonster(worldId byte, channelId byte, mapId uint32, monsterId int, x int, y int, fh int, stance int, team int, hp uint32, mp uint32) *Model {

	r.monsterRegisterRWLock.Lock()
	var existingIds = existingIds(r.monsterRegister)

	var currentUniqueId = uniqueId
	for contains(existingIds, currentUniqueId) {
		currentUniqueId = currentUniqueId + 1
		if currentUniqueId > 2000000000 {
			currentUniqueId = 1000000001
		}
		uniqueId = currentUniqueId
	}

	m := NewMonster(worldId, channelId, mapId, currentUniqueId, monsterId, x, y, fh, stance, team, hp, mp)

	r.monsterRegister[uniqueId] = m
	r.monsterRegisterRWLock.Unlock()

	mk := models.NewMapKey(worldId, channelId, mapId)
	r.getMapLock(*mk).Lock()
	if om, ok := r.mapMonsters[*mk]; ok {
		r.mapMonsters[*mk] = append(om, m.UniqueId())
	} else {
		r.mapMonsters[*mk] = append([]int{}, m.UniqueId())
	}
	r.getMapLock(*mk).Unlock()

	return m
}

func (r *Registry) GetMonster(uniqueId int) (*Model, error) {
	r.monsterRegisterRWLock.RLock()
	if val, ok := r.monsterRegister[uniqueId]; ok {
		r.monsterRegisterRWLock.RUnlock()
		return val, nil
	} else {
		r.monsterRegisterRWLock.RUnlock()
		return nil, errors.New("monster not found")
	}
}

func (r *Registry) GetMonstersInMap(worldId byte, channelId byte, mapId uint32) []*Model {
	mk := models.NewMapKey(worldId, channelId, mapId)
	r.getMapLock(*mk).Lock()
	r.monsterRegisterRWLock.RLock()
	var result []*Model
	for _, x := range r.mapMonsters[*mk] {
		result = append(result, r.monsterRegister[x])
	}
	r.monsterRegisterRWLock.RUnlock()
	r.getMapLock(*mk).Unlock()
	return result
}

func (r *Registry) MoveMonster(uniqueId int, endX int, endY int, stance int) {
	r.monsterRegisterRWLock.Lock()
	if m, ok := r.monsterRegister[uniqueId]; ok {
		um := m.Move(endX, endY, stance)
		r.monsterRegister[uniqueId] = um
	}
	r.monsterRegisterRWLock.Unlock()
}

func (r *Registry) ControlMonster(uniqueId int, characterId uint32) *Model {
	r.monsterRegisterRWLock.Lock()
	if m, ok := r.monsterRegister[uniqueId]; ok {
		um := m.Control(characterId)
		r.monsterRegister[uniqueId] = um
		r.monsterRegisterRWLock.Unlock()
		return um
	} else {
		r.monsterRegisterRWLock.Unlock()
		return nil
	}
}

func (r *Registry) ClearControl(uniqueId int) *Model {
	r.monsterRegisterRWLock.Lock()
	if m, ok := r.monsterRegister[uniqueId]; ok {
		um := m.ClearControl()
		r.monsterRegister[uniqueId] = um
		r.monsterRegisterRWLock.Unlock()
		return um
	} else {
		r.monsterRegisterRWLock.Unlock()
		return nil
	}
}

func (r *Registry) ApplyDamage(characterId uint32, damage int64, uniqueId int) (*DamageSummary, error) {
	r.monsterRegisterRWLock.Lock()
	if m, ok := r.monsterRegister[uniqueId]; ok {
		um := m.Damage(characterId, damage)
		r.monsterRegister[uniqueId] = um
		r.monsterRegisterRWLock.Unlock()
		return &DamageSummary{
			CharacterId:   characterId,
			Monster:       *um,
			VisibleDamage: damage,
			ActualDamage:  int64(m.Hp() - um.Hp()),
			Killed:        um.Hp() == 0,
		}, nil
	} else {
		r.monsterRegisterRWLock.Unlock()
		return nil, errors.New("monster not found")
	}
}

func (r *Registry) RemoveMonster(uniqueId int) {
	r.monsterRegisterRWLock.Lock()
	if m, ok := r.monsterRegister[uniqueId]; ok {
		mk := models.NewMapKey(m.WorldId(), m.ChannelId(), m.MapId())
		r.removeMonster(*mk, uniqueId)
	}
	r.monsterRegisterRWLock.Unlock()
}

func remove(c []int, i int) []int {
	c[i] = c[len(c)-1]
	return c[:len(c)-1]
}

func indexOf(id int, data []int) int {
	for k, v := range data {
		if id == v {
			return k
		}
	}
	return -1 //not found.
}

func (r *Registry) removeMonster(mapId models.MapKey, uniqueId int) {
	index := indexOf(uniqueId, r.mapMonsters[mapId])
	if index >= 0 && index < len(r.mapMonsters[mapId]) {
		r.mapMonsters[mapId] = remove(r.mapMonsters[mapId], index)
	}
}

func (r *Registry) GetMonsters() []*Model {
	r.monsterRegisterRWLock.RLock()
	ms := make([]*Model, len(r.monsterRegister))
	for _, x := range r.monsterRegister {
		ms = append(ms, x)
	}
	r.monsterRegisterRWLock.RUnlock()
	return ms
}
