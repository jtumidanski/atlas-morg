package monster

type DamageSummary struct {
	CharacterId   uint32
	Monster       Model
	VisibleDamage int64
	ActualDamage  int64
	Killed        bool
}
