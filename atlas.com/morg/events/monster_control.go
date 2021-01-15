package events

type MonsterControlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId int    `json:"characterId"`
	UniqueId    int    `json:"uniqueId"`
	Type        string `json:"type"`
}
