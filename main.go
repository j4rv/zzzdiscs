package zzzdiscs

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

const MaxSubstats = 4
const Base4Chance = 1.0 / 5.0 // 20% for CleanUp and Elfy
const AverageDropsPerRoutineRun = 2.47
const DomainExtraDiscChance = 0.47

type DiscStat struct {
	Stat      stat
	BaseValue int
	Value     int
	Rolls     int
}

func (s *DiscStat) String() string {
	suffix := ""
	if s.Rolls > 1 {
		suffix = " +" + strconv.Itoa(s.Rolls-1)
	}
	return s.Stat.String() + suffix
}

type Disc struct {
	Set         DiscSet
	Slot        DiscSlot
	MainStat    DiscStat
	SubStats    [MaxSubstats]*DiscStat
	IsFourLiner bool
}

func (a Disc) String() string {
	subsStr := ""
	for _, s := range a.SubStats {
		subsStr += s.String() + "\n"
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s", a.Set, a.MainStat.String(), subsStr)
}

func (a Disc) hasSubstat(s stat) bool {
	for _, sub := range a.SubStats {
		if sub.Stat == s {
			return true
		}
	}
	return false
}

func (a Disc) subsQuality(wantedSubWeights map[stat]float32) float32 {
	var quality float32
	for _, sub := range a.SubStats {
		quality += wantedSubWeights[sub.Stat]
	}
	return quality
}

func (a Disc) cv() int {
	var cv int
	for _, sub := range a.SubStats {
		switch sub.Stat {
		case CritRate:
			cv += sub.Value * 2
		case CritDmg:
			cv += sub.Value
		}
	}
	return cv
}

func (a *Disc) randomizeSet(options ...DiscSet) {
	a.Set = options[rand.Intn(len(options))]
}

func (a *Disc) randomizeSlot() {
	a.Slot = DiscSlot(rand.Intn(5))
}

func (a *Disc) ranzomizeMainStat() {
	switch a.Slot {
	case Slot1:
		a.MainStat.Stat = HP
	case Slot2:
		a.MainStat.Stat = ATK
	case Slot3:
		a.MainStat.Stat = DEF
	case Slot4:
		a.MainStat.Stat = d4StatsRandChooser.Pick()
	case Slot5:
		a.MainStat.Stat = d5StatsRandChooser.Pick()
	case Slot6:
		a.MainStat.Stat = d6StatsRandChooser.Pick()
	}
	a.MainStat.BaseValue = mainStatValues[a.MainStat.Stat]
	a.MainStat.Value = mainStatValues[a.MainStat.Stat]
}

func (a *Disc) randomizeSubstats(base4Chance float32, maxLevel bool) {
	numRolls := 3 + 5 // starts with 3 subs by default
	if rand.Float32() <= base4Chance {
		numRolls++ // starts with 4 subs
		a.IsFourLiner = true
	}

	a.SubStats = [MaxSubstats]*DiscStat{}
	possibleStats := weightedSubstats(a.MainStat.Stat)
	var subs [MaxSubstats]stat

	for i := 0; i < numRolls; i++ {
		// Stop before adding the 4th substat if it's a 3-liner and not max level
		if i == 3 && !maxLevel && !a.IsFourLiner {
			break
		}
		// First 3-4 rolls
		if i < MaxSubstats {
			artiStat := weightedRand(possibleStats)
			subs[i] = artiStat
			a.SubStats[i] = &DiscStat{Stat: artiStat, Rolls: 1, BaseValue: substatValues[artiStat], Value: substatValues[artiStat]}
			delete(possibleStats, artiStat)
		} else {
			// Rest of rolls, if max level
			index := rand.Intn(MaxSubstats)
			a.SubStats[index].Rolls += 1
			a.SubStats[index].Value += a.SubStats[index].BaseValue
		}
	}
}

func RandomDisc(base4Chance float32) *Disc {
	var disc Disc
	disc.randomizeSet(AllDiscSets...)
	disc.randomizeSlot()
	disc.ranzomizeMainStat()
	disc.randomizeSubstats(base4Chance, true)
	return &disc
}

func RandomDiscOfSetAndSlot(set DiscSet, slot DiscSlot, base4Chance float32) *Disc {
	var disc Disc
	disc.Set = set
	disc.Slot = slot
	disc.ranzomizeMainStat()
	disc.randomizeSubstats(base4Chance, true)
	return &disc
}

func RandomDiscOfSlot(slot DiscSlot, base4Chance float32) *Disc {
	var disc Disc
	disc.randomizeSet(AllDiscSets...)
	disc.Slot = slot
	disc.ranzomizeMainStat()
	disc.randomizeSubstats(base4Chance, true)
	return &disc
}

func RandomDiscOfSet(set string, base4Chance float32) *Disc {
	var disc Disc
	disc.Set = DiscSet(set)
	disc.randomizeSlot()
	disc.ranzomizeMainStat()
	disc.randomizeSubstats(base4Chance, true)
	return &disc
}

func RandomDiscFromDomain(setA, setB string) *Disc {
	var disc Disc
	disc.randomizeSet(DiscSet(setA), DiscSet(setB))
	disc.randomizeSlot()
	disc.ranzomizeMainStat()
	disc.randomizeSubstats(Base4Chance, true)
	return &disc
}

// RemoveTrashDiscs processes a slice of discs and keeps the best ones that have the correct mainstat
// subValue: To know which discs are more desirable
// n: Amount of discs to keep for every set, slot and main stat (example: n = 10, it will keep at most 10 gladiator atk sands)
func RemoveTrashDiscs(arts []*Disc,
	subValue map[stat]float32,
	n int) []*Disc {
	type SetSlotStat struct {
		set      DiscSet
		slot     DiscSlot
		mainStat stat
	}
	processed := map[SetSlotStat][]*Disc{}
	for _, art := range arts {
		sss := SetSlotStat{art.Set, art.Slot, art.MainStat.Stat}
		processed[sss] = append(processed[sss], art)
	}

	result := []*Disc{}
	for _, aa := range processed {
		// Ordering the discs in processed by sub quality
		sort.Slice(aa, func(i, j int) bool {
			return aa[i].subsQuality(subValue) > aa[j].subsQuality(subValue)
		})
		// Keeping the n best
		if len(aa) > n {
			aa = aa[0:n]
		}
		result = append(result, aa...)
	}
	return result
}
