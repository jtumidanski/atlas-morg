package events

type MonsterEvent struct {
   WorldId   byte   `json:"worldId"`
   ChannelId byte   `json:"channelId"`
   MapId     int    `json:"mapId"`
   UniqueId  int    `json:"uniqueId"`
   MonsterId int    `json:"monsterId"`
   ActorId   int    `json:"actorId"`
   Type      string `json:"type"`
}
