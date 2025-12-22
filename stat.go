package zzzdiscs

import (
	"github.com/mroth/weightedrand/v2"
)

type stat int

const (
	HP stat = iota
	ATK
	DEF
	PEN
	HPP
	ATKP
	DEFP
	PENP
	CritRate
	CritDmg
	FireDMG
	ElecDMG
	IceDMG
	PhysDMG
	EtherDMG
	Impact
	EnergyRecharge
	AnomalyMastery
	AnomalyProficiency

	GlobalDMGBonus
	BaseDMGIncrease
)

var StatToID = map[stat]int{
	HP:                 11103,
	ATK:                12103,
	DEF:                13103,
	HPP:                11102,
	ATKP:               12102,
	DEFP:               13102,
	PEN:                23203,
	PENP:               23103,
	CritRate:           20103,
	CritDmg:            21103,
	FireDMG:            31603,
	ElecDMG:            31803,
	IceDMG:             31703,
	PhysDMG:            31503,
	EtherDMG:           31903,
	Impact:             12202,
	EnergyRecharge:     30502,
	AnomalyMastery:     31402,
	AnomalyProficiency: 31203,
}

var substatValues = map[stat]int{
	HP:                 112,
	HPP:                300,
	ATK:                19,
	ATKP:               300,
	DEF:                15,
	DEFP:               480,
	PEN:                9,
	AnomalyProficiency: 9,
	CritDmg:            480,
	CritRate:           240,
}

var mainStatValues = map[stat]int{
	HP:                 550,
	HPP:                750,
	ATK:                79,
	ATKP:               750,
	DEF:                46,
	DEFP:               1200,
	PEN:                23,
	PENP:               600,
	CritDmg:            1200,
	CritRate:           600,
	EnergyRecharge:     1500,
	Impact:             450,
	ElecDMG:            750,
	EtherDMG:           750,
	FireDMG:            750,
	IceDMG:             750,
	PhysDMG:            750,
	AnomalyMastery:     750,
	AnomalyProficiency: 23,
}

func (s stat) String() string {
	switch s {
	case HP:
		return "HP"
	case ATK:
		return "ATK"
	case DEF:
		return "DEF"
	case HPP:
		return "HP%"
	case ATKP:
		return "ATK%"
	case DEFP:
		return "DEF%"
	case PEN:
		return "PEN"
	case PENP:
		return "PEN Ratio"
	case CritRate:
		return "Crit Rate"
	case CritDmg:
		return "Crit DMG"
	case FireDMG:
		return "Fire DMG"
	case ElecDMG:
		return "Electric DMG"
	case IceDMG:
		return "Ice DMG"
	case PhysDMG:
		return "Physical DMG"
	case EtherDMG:
		return "Ether DMG"
	case Impact:
		return "Impact"
	case EnergyRecharge:
		return "Energy Regen%"
	case AnomalyMastery:
		return "Anomaly Mastery"
	case AnomalyProficiency:
		return "Anomaly Proficiency"
	case GlobalDMGBonus:
		return "Global DMG Bonus"
	case BaseDMGIncrease:
		return "Base DMG Increase"
	default:
		return "Unknown"
	}
}

// Weights from https://docs.google.com/spreadsheets/d/1oeSKlAHYqFqHvwh866KElWFYK3ZaYH7rGIsLS_h18qg
// And my own weights study
// Left: Sixth Alley data, Right: My data

var d4StatsRandChooser, _ = weightedrand.NewChooser(
	weightedrand.NewChoice(stat(HPP), 1740+100),
	weightedrand.NewChoice(stat(ATKP), 1479+90),
	weightedrand.NewChoice(stat(DEFP), 1660+116),
	weightedrand.NewChoice(stat(AnomalyProficiency), 1216+69),
	weightedrand.NewChoice(stat(CritRate), 972+58),
	weightedrand.NewChoice(stat(CritDmg), 992+62),
)

var d5DmgWeight int = ((523 + 497 + 493 + 492 + 486) + (29 + 38 + 25 + 32 + 26)) / 5
var d5StatsRandChooser, _ = weightedrand.NewChooser(
	weightedrand.NewChoice(stat(HPP), 1710+106),
	weightedrand.NewChoice(stat(ATKP), 1447+80),
	weightedrand.NewChoice(stat(DEFP), 1733+92),
	weightedrand.NewChoice(stat(PENP), 797+43),
	weightedrand.NewChoice(stat(FireDMG), d5DmgWeight),
	weightedrand.NewChoice(stat(ElecDMG), d5DmgWeight),
	weightedrand.NewChoice(stat(PhysDMG), d5DmgWeight),
	weightedrand.NewChoice(stat(EtherDMG), d5DmgWeight),
	weightedrand.NewChoice(stat(IceDMG), d5DmgWeight),
)

var d6StatsRandChooser, _ = weightedrand.NewChooser(
	weightedrand.NewChoice(stat(HPP), 1659+99),
	weightedrand.NewChoice(stat(ATKP), 1393+95),
	weightedrand.NewChoice(stat(DEFP), 1609+105),
	weightedrand.NewChoice(stat(Impact), 1191+63),
	weightedrand.NewChoice(stat(EnergyRecharge), 773+44),
	weightedrand.NewChoice(stat(AnomalyMastery), 1162+67),
)

func weightedSubstats(mainStat stat) map[stat]int {
	weightedSubs := map[stat]int{
		HP:                 2246 + 106,
		ATK:                2015 + 90,
		DEF:                2262 + 83,
		HPP:                2241 + 95,
		ATKP:               2012 + 86,
		DEFP:               2171 + 86,
		PEN:                1954 + 83,
		CritRate:           1896 + 72,
		CritDmg:            1860 + 61,
		AnomalyProficiency: 1857 + 70,
	}
	delete(weightedSubs, mainStat)
	return weightedSubs
}
