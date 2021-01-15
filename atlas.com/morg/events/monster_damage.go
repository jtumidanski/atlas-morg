package events

type MonsterDamageEvent struct {
   WorldId     byte  `json:"worldId"`
   ChannelId   byte  `json:"channelId"`
   MapId       int   `json:"mapId"`
   UniqueId    int   `json:"uniqueId"`
   CharacterId int   `json:"characterId"`
   Damage      int64 `json:"damage"`
}
