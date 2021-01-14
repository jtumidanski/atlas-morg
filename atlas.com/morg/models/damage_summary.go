package models

type DamageSummary struct {
	CharacterId   int
	Monster       Monster
	VisibleDamage int32
	ActualDamage  int32
	Killed        bool
}
