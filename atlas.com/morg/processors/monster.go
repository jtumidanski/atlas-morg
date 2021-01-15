package processors

import (
   "atlas-morg/attributes"
   "atlas-morg/models"
   "atlas-morg/producers"
   "atlas-morg/registries"
   "atlas-morg/requests"
   "context"
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
   ma, err := requests.NewMapInformation(m.l).GetMonsterInformation(input.MonsterId)
   if err != nil {
      return nil, err
   }

   model := registries.GetMonsterRegistry().CreateMonster(worldId, channelId, mapId, input.MonsterId, input.X, input.Y, input.Fh, 5, input.Team, ma.Hp)

   cid, err := m.GetControllerCandidate(worldId, channelId, mapId)
   if err == nil {
      m.StartControl(model, cid)
   }

   producers.NewMonster(m.l, context.Background()).EmitCreated(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
   return model, nil
}

func (m *Monster) FindNextController(model *models.Monster) {
   cid, err := m.GetControllerCandidate(model.WorldId(), model.ChannelId(), model.MapId())
   if err == nil {
      m.StartControl(model, cid)
   }
}

func (m *Monster) GetControllerCandidate(worldId byte, channelId byte, mapId int) (int, error) {
   ids, err := requests.NewMap(m.l).GetCharacterIdsInMap(worldId, channelId, mapId)
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
   producers.NewMonsterControl(m.l, context.Background()).EmitControl(monster.WorldId(), monster.ChannelId(), monster.MapId(), monster.ControlCharacterId(), monster.UniqueId())
}

func (m *Monster) StopControl(monster *models.Monster) {
   registries.GetMonsterRegistry().ClearControl(monster.UniqueId())
   producers.NewMonsterControl(m.l, context.Background()).EmitStop(monster.WorldId(), monster.ChannelId(), monster.MapId(), monster.ControlCharacterId(), monster.UniqueId())
}

func (m *Monster) DestroyAll() {
   ms := registries.GetMonsterRegistry().GetMonsters()
   for _, x := range ms {
      m.Destroy(&x)
   }
}

func (m *Monster) Destroy(model *models.Monster) {
   registries.GetMonsterRegistry().RemoveMonster(model.UniqueId())
   producers.NewMonster(m.l, context.Background()).EmitDestroyed(model.WorldId(), model.ChannelId(), model.MapId(), model.UniqueId(), model.MonsterId())
}
