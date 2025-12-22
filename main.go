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
		quality += wantedSubWeights[sub.Stat] * float32(sub.Rolls)
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

func (a *Disc) setMainStat(s stat) {
	a.MainStat.Stat = s
	a.MainStat.BaseValue = mainStatValues[s]
	a.MainStat.Value = mainStatValues[s]
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
		a.setMainStat(HP)
	case Slot2:
		a.setMainStat(ATK)
	case Slot3:
		a.setMainStat(DEF)
	case Slot4:
		a.setMainStat(d4StatsRandChooser.Pick())
	case Slot5:
		a.setMainStat(d5StatsRandChooser.Pick())
	case Slot6:
		a.setMainStat(d6StatsRandChooser.Pick())
	}
}

func (a *Disc) randomizeSubstats(base4Chance float32, maxLevel bool, chosen ...stat) {
	numRolls := 3 + 5
	if rand.Float32() <= base4Chance {
		numRolls++ // starts with 4 subs
		a.IsFourLiner = true
	}

	a.SubStats = [MaxSubstats]*DiscStat{}
	possibleStats := weightedSubstats(a.MainStat.Stat)

	if len(chosen) > MaxSubstats {
		chosen = chosen[:MaxSubstats]
	}

	// First add the chosen substats
	for i, s := range chosen {
		a.SubStats[i] = &DiscStat{
			Stat:      s,
			Rolls:     1,
			BaseValue: substatValues[s],
			Value:     substatValues[s],
		}
		delete(possibleStats, s)
	}

	for i := len(chosen); i < numRolls; i++ {
		// Stop before adding the 4th substat if it's a 3-liner and not max level
		if i == 3 && !maxLevel && !a.IsFourLiner {
			break
		}

		if i < MaxSubstats {
			// Add new substat
			artiStat := weightedRand(possibleStats)
			a.SubStats[i] = &DiscStat{
				Stat:      artiStat,
				Rolls:     1,
				BaseValue: substatValues[artiStat],
				Value:     substatValues[artiStat],
			}
			delete(possibleStats, artiStat)
		} else {
			// Roll existing substat
			index := rand.Intn(MaxSubstats)
			a.SubStats[index].Rolls++
			a.SubStats[index].Value += a.SubStats[index].BaseValue
		}
	}
}

// Deprecated
func RandomDiscSimple() *Disc {
	var disc Disc
	disc.randomizeSet(AllDiscSets...)
	disc.randomizeSlot()
	disc.ranzomizeMainStat()
	disc.randomizeSubstats(Base4Chance, true)
	return &disc
}

// Deprecated
func RandomDiscOfSetAndSlot(set DiscSet, slot DiscSlot, base4Chance float32) *Disc {
	return RandomDisc()
}

// Deprecated
func RandomDiscOfSlot(slot DiscSlot, base4Chance float32) *Disc {
	return RandomDisc(
		WithSlot(slot),
		WithBase4Chance(base4Chance),
	)
}

// Deprecated
func RandomDiscOfSet(set string, base4Chance float32) *Disc {
	return RandomDisc(
		WithSet(DiscSet(set)),
		WithBase4Chance(base4Chance),
	)
}

// Deprecated
func RandomDiscFromDomain(setA, setB string) *Disc {
	disc := RandomDisc()
	disc.randomizeSet(DiscSet(setA), DiscSet(setB))
	return disc
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
