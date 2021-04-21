package monster

type DamageSummary struct {
	CharacterId   int
	Monster       Model
	VisibleDamage int64
	ActualDamage  int64
	Killed        bool
}
