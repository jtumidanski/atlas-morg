package producers

import (
   "atlas-morg/events"
   "atlas-morg/requests"
   "context"
   "encoding/json"
   "github.com/segmentio/kafka-go"
   "log"
   "os"
   "strconv"
   "time"
)

type MonsterControl struct {
   l   *log.Logger
   ctx context.Context
}

func NewMonsterControl(l *log.Logger, ctx context.Context) *MonsterControl {
   return &MonsterControl{l, ctx}
}

func (m *MonsterControl) EmitControl(worldId byte, channelId byte, mapId int, characterId int, uniqueId int) {
   m.emit(worldId, channelId, mapId, characterId, uniqueId, "START")
}

func (m *MonsterControl) EmitStop(worldId byte, channelId byte, mapId int, characterId int, uniqueId int) {
   m.emit(worldId, channelId, mapId, characterId, uniqueId, "STOP")
}

func (m *MonsterControl) emit(worldId byte, channelId byte, mapId int, characterId int, uniqueId int, theType string) {
   t := requests.NewTopic(m.l)
   td, err := t.GetTopic("TOPIC_CONTROL_MONSTER_EVENT")
   if err != nil {
      m.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
   }

   w := &kafka.Writer{
      Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
      Topic:        td.Attributes.Name,
      Balancer:     &kafka.LeastBytes{},
      BatchTimeout: 50 * time.Millisecond,
   }

   e := &events.MonsterControlEvent{
      WorldId:     worldId,
      ChannelId:   channelId,
      MapId:       mapId,
      CharacterId: characterId,
      UniqueId:    uniqueId,
      Type:        theType,
   }
   r, err := json.Marshal(e)
   if err != nil {
      m.l.Fatal("[ERROR] Unable to marshall event.")
   }

   m.l.Printf("[INFO] Sending [MonsterControlEvent] key %s", createKey(mapId))
   m.l.Printf("[INFO] Sending [MonsterControlEvent] value %s", r)

   err = w.WriteMessages(context.Background(), kafka.Message{
      Key:   createKey(mapId),
      Value: r,
   })
   if err != nil {
      m.l.Fatal("[ERROR] Unable to produce event.")
   }
}

func createKey(key int) []byte {
   var empty = make([]byte, 8)
   sk := []byte(strconv.Itoa(key))

   start := len(empty) - len(sk)
   var result = empty

   for i := 0; i < start; i++ {
      result[i] = 48
   }

   for i := start; i < len(empty); i++ {
      result[i] = sk[i-start]
   }
   return result
}
