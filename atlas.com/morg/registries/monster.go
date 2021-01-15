package registries

import (
   "atlas-morg/models"
   "errors"
   "sync"
)

type MonsterRegistry struct {
   mutex sync.Mutex

   monsterMap  map[int]models.Monster
   mapMonsters map[models.MapKey][]int

   monsterLocks map[int]*sync.Mutex
}

var monsterRegistry *MonsterRegistry
var once sync.Once

var uniqueId = 1000000001

func GetMonsterRegistry() *MonsterRegistry {
   once.Do(func() {
      monsterRegistry = &MonsterRegistry{}
      monsterRegistry.monsterMap = make(map[int]models.Monster)
      monsterRegistry.mapMonsters = make(map[models.MapKey][]int)
      monsterRegistry.monsterLocks = make(map[int]*sync.Mutex)
   })
   return monsterRegistry
}

func (r *MonsterRegistry) getMonsterLock(uniqueId int) *sync.Mutex {
   if val, ok := r.monsterLocks[uniqueId]; ok {
      return val
   } else {
      var cm = &sync.Mutex{}
      r.mutex.Lock()
      r.monsterLocks[uniqueId] = cm
      r.mutex.Unlock()
      return cm
   }
}

func existingIds(monsters map[int]models.Monster) []int {
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

func (r *MonsterRegistry) CreateMonster(worldId byte, channelId byte, mapId int, monsterId int, x int, y int, fh int, stance int, team int, hp int) *models.Monster {
   r.mutex.Lock()

   var existingIds = existingIds(r.monsterMap)

   var currentUniqueId = uniqueId
   for contains(existingIds, currentUniqueId) {
      currentUniqueId = currentUniqueId + 1
      if currentUniqueId > 2000000000 {
         currentUniqueId = 1000000001
      }
      uniqueId = currentUniqueId
   }

   m := models.NewMonster(worldId, channelId, mapId, currentUniqueId, monsterId, x, y, fh, stance, team, hp)
   r.monsterMap[uniqueId] = *m

   mk := models.NewMapKey(worldId, channelId, mapId)
   if om, ok := r.mapMonsters[*mk]; ok {
      r.mapMonsters[*mk] = append(om, m.UniqueId())
   } else {
      r.mapMonsters[*mk] = append([]int{}, m.UniqueId())
   }

   r.mutex.Unlock()
   return m
}

func (r *MonsterRegistry) GetMonster(uniqueId int) (*models.Monster, error) {
   if val, ok := r.monsterMap[uniqueId]; ok {
      return &val, nil
   } else {
      return nil, errors.New("monster not found")
   }
}

func (r *MonsterRegistry) GetMonstersInMap(worldId byte, channelId byte, mapId int) []models.Monster {
   mk := models.NewMapKey(worldId, channelId, mapId)
   var result []models.Monster
   for _, x := range r.mapMonsters[*mk] {
      result = append(result, r.monsterMap[x])
   }
   return result
}

func (r *MonsterRegistry) MoveMonster(uniqueId int, endX int, endY int, stance int) {
   ml := r.getMonsterLock(uniqueId)
   ml.Lock()
   if m, ok := r.monsterMap[uniqueId]; ok {
      um := m.Move(endX, endY, stance)
      r.monsterMap[uniqueId] = *um
   }
   ml.Unlock()
}

func (r *MonsterRegistry) ControlMonster(uniqueId int, characterId int) {
   ml := r.getMonsterLock(uniqueId)
   ml.Lock()
   if m, ok := r.monsterMap[uniqueId]; ok {
      um := m.Control(characterId)
      r.monsterMap[uniqueId] = *um
   }
   ml.Unlock()
}

func (r *MonsterRegistry) ClearControl(uniqueId int) {
   ml := r.getMonsterLock(uniqueId)
   ml.Lock()
   if m, ok := r.monsterMap[uniqueId]; ok {
      um := m.ClearControl()
      r.monsterMap[uniqueId] = *um
   }
   ml.Unlock()
}

func (r *MonsterRegistry) ApplyDamage(characterId int, damage int64, uniqueId int) (*models.DamageSummary, error) {
   ml := r.getMonsterLock(uniqueId)
   ml.Lock()

   if m, ok := r.monsterMap[uniqueId]; ok {
      um := m.Damage(characterId, damage)
      r.monsterMap[uniqueId] = *um
      ml.Unlock()
      return &models.DamageSummary{
         CharacterId:   characterId,
         Monster:       *um,
         VisibleDamage: damage,
         ActualDamage:  int64(m.Hp() - um.Hp()),
         Killed:        um.Hp() == 0,
      }, nil
   } else {
      ml.Unlock()
      return nil, errors.New("monster not found")
   }
}

func (r *MonsterRegistry) RemoveMonster(uniqueId int) {
   ml := r.getMonsterLock(uniqueId)
   ml.Lock()
   if m, ok := r.monsterMap[uniqueId]; ok {
      mk := models.NewMapKey(m.WorldId(), m.ChannelId(), m.MapId())
      r.removeMonster(*mk, uniqueId)
   }
   ml.Unlock()
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

func (r *MonsterRegistry) removeMonster(mapId models.MapKey, uniqueId int) {
   index := indexOf(uniqueId, r.mapMonsters[mapId])
   if index >= 0 && index < len(r.mapMonsters[mapId]) {
      r.mapMonsters[mapId] = remove(r.mapMonsters[mapId], index)
   }
}

func (r *MonsterRegistry) GetMonsters() []models.Monster {
   ms := make([]models.Monster, len(r.monsterMap))
   for _, x := range r.monsterMap {
      ms = append(ms, x)
   }
   return ms
}
