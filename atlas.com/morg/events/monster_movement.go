package events

type MonsterMovementEvent struct {
	UniqueId      int    `json:"uniqueId"`
	ObserverId    int    `json:"observerId"`
	SkillPossible bool   `json:"skillPossible"`
	Skill         int    `json:"skill"`
	SkillId       int    `json:"skillId"`
	SkillLevel    int    `json:"skillLevel"`
	Option        int    `json:"option"`
	StartX        int    `json:"startX"`
	StartY        int    `json:"startY"`
	EndX          int    `json:"endX"`
	EndY          int    `json:"endY"`
	Stance        int    `json:"stance"`
	RawMovement   []int `json:"rawMovement"`
}
