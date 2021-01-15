package events

type DamageEntry struct {
   CharacterId int   `json:"characterId"`
   Damage      int64 `json:"damage"`
}
