package attributes

type MonsterInputDataContainer struct {
	Data MonsterData `json:"data"`
}

type MonsterDataContainer struct {
	Data MonsterData `json:"data"`
}

type MonsterListDataContainer struct {
	Data []MonsterData `json:"data"`
}

type MonsterData struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes MonsterAttributes `json:"attributes"`
}

type MonsterAttributes struct {
	WorldId            byte            `json:"worldId"`
	ChannelId          byte            `json:"channelId"`
	MapId              uint32          `json:"mapId"`
	MonsterId          int             `json:"monsterId"`
	ControlCharacterId uint32          `json:"controlCharacterId"`
	X                  int             `json:"x"`
	Y                  int             `json:"y"`
	Fh                 int             `json:"fh"`
	Stance             int             `json:"stance"`
	Team               int             `json:"team"`
	MaxHp              uint32          `json:"maxHp"`
	Hp                 uint32          `json:"hp"`
	MaxMp              uint32          `json:"maxMp"`
	Mp                 uint32          `json:"mp"`
	DamageEntries      []MonsterDamage `json:"damageEntries"`
}

type MonsterDamage struct {
	CharacterId uint32 `json:"characterId"`
	Damage      int64  `json:"damage"`
}
