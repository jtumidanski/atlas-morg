package events

type MonsterKilledEvent struct {
   WorldId       byte          `json:"worldId"`
   ChannelId     byte          `json:"channelId"`
   MapId         int           `json:"mapId"`
   UniqueId      int           `json:"uniqueId"`
   MonsterId     int           `json:"monsterId"`
   X             int           `json:"x"`
   Y             int           `json:"y"`
   KillerId      int           `json:"killerId"`
   DamageEntries []DamageEntry `json:"damageEntries"`
}
