package zzzdiscs

type RandomDiscOption func(*randomDiscOptions)

type randomDiscOptions struct {
	set            *DiscSet
	slot           *DiscSlot
	mainStat       *stat
	chosenSubStats []stat
	base4Chance    float32
}

func WithSet(set DiscSet) RandomDiscOption {
	return func(o *randomDiscOptions) {
		o.set = &set
	}
}

func WithSlot(slot DiscSlot) RandomDiscOption {
	return func(o *randomDiscOptions) {
		o.slot = &slot
	}
}

func WithMainStat(stat stat) RandomDiscOption {
	return func(o *randomDiscOptions) {
		o.mainStat = &stat
	}
}

func WithSubstats(stats ...stat) RandomDiscOption {
	return func(o *randomDiscOptions) {
		o.chosenSubStats = stats
	}
}

func WithBase4Chance(chance float32) RandomDiscOption {
	return func(o *randomDiscOptions) {
		o.base4Chance = chance
	}
}

func RandomDisc(opts ...RandomDiscOption) *Disc {
	// defaults
	cfg := randomDiscOptions{
		base4Chance: Base4Chance,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	var disc Disc

	// Set
	if cfg.set != nil {
		disc.Set = *cfg.set
	} else {
		disc.randomizeSet(AllDiscSets...)
	}

	// Slot
	if cfg.slot != nil {
		disc.Slot = *cfg.slot
	} else {
		disc.randomizeSlot()
	}

	// Main stat
	if cfg.mainStat != nil {
		disc.setMainStat(*cfg.mainStat)
	} else {
		disc.ranzomizeMainStat()
	}

	// Substats
	disc.randomizeSubstats(cfg.base4Chance, true, cfg.chosenSubStats...)

	return &disc
}
