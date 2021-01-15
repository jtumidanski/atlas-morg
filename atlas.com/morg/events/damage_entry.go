package events

type DamageEntry struct {
	Character int   `json:"character"`
	Damage    int64 `json:"damage"`
}
