package models

type DamageSummary struct {
	CharacterId   int
	Monster       Monster
	VisibleDamage int64
	ActualDamage  int64
	Killed        bool
}
