package zzzdiscs

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestRandomDiscFromDomain(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var set1Count, set2Count int
	set1, set2 := "Emblem", "Shimenawa"

	// Generate 1000 discs from two sets
	for i := 0; i < 1000; i++ {
		art := RandomDiscFromDomain(set1, set2)
		if art.Set == DiscSet(set1) {
			set1Count++
		} else if art.Set == DiscSet(set2) {
			set2Count++
		} else {
			t.Error("Unexpected disc set: " + art.Set)
		}
	}

	// Then check that the chances of getting an disc from either set is ~50%
	if set1Count < 450 || set1Count > 550 {
		t.Error("Too many or too few discs from set 1: " + strconv.Itoa(set1Count))
	}
	if set2Count < 450 || set2Count > 550 {
		t.Error("Too many or too few discs from set 1: " + strconv.Itoa(set2Count))
	}
}

func TestRemoveTrashDiscs(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var discs []*Disc
	set1, set2 := "Emblem", "Shimenawa"

	// Generate 1000 discs from two sets
	for i := 0; i < 10000; i++ {
		discs = append(discs, RandomDiscFromDomain(set1, set2))
	}

	subs := map[stat]float32{
		HP:                 0.1,
		DEF:                0,
		ATK:                0.1,
		PEN:                0.4,
		HPP:                0.8,
		ATKP:               0.8,
		DEFP:               0.2,
		AnomalyProficiency: 1,
		CritRate:           1,
		CritDmg:            1,
	}
	filtered := RemoveTrashDiscs(discs, subs, 5)
	for _, a := range filtered {
		t.Log(*a)
	}
}

func TestExportToRR(t *testing.T) {
	fileIndex := 0
	keepPerSetSlotAndMainStat := 5
	for _, set := range AllDiscSets {
		for slot := DiscSlot(0); slot < 6; slot++ {
			var discs []*Disc
			for j := 0; j < 100; j++ {
				discs = append(discs, RandomDiscOfSetAndSlot(set, slot, Base4Chance))
			}

			subs := map[stat]float32{
				HP:                 0.2,
				DEF:                0,
				ATK:                0.2,
				PEN:                0.4,
				HPP:                0.8,
				ATKP:               0.8,
				DEFP:               0,
				AnomalyProficiency: 1,
				CritRate:           1,
				CritDmg:            1,
			}
			discs = RemoveTrashDiscs(discs, subs, keepPerSetSlotAndMainStat)

			for _, disc := range discs {
				export := FormatRRDisc(DiscToRRDisc(*disc))
				filename := "./discs/" + strconv.Itoa(10000+fileIndex)
				os.WriteFile(filename, []byte(export), 0755)
				fileIndex++
				//
			}
		}
	}
}

func TestWithChosenSubstats(t *testing.T) {
	disc := RandomDisc(
		WithMainStat(stat(PENP)),
		WithSubstats(
			CritRate,
			CritDmg,
		),
	)
	log.Println(disc.String())
}

func WithChosenSubstatsValueSingleTest(subQualityMap map[stat]float32, totalDiscs int, wantedMainstat stat, wantedSubstats ...stat) {
	qualityCounts := map[int]int{
		5: 0,
		6: 0,
		7: 0,
	}
	tunerCost := totalDiscs * (1 + 2*len(wantedSubstats))

	for i := 0; i < totalDiscs; i++ {
		disc := RandomDisc(
			WithMainStat(wantedMainstat),
			WithSubstats(wantedSubstats...),
		)
		//log.Println(disc.String())
		//log.Printf("Quality: %.2v", disc.subsQuality(subQualityMap))
		quality := disc.subsQuality(subQualityMap)
		if quality >= 5 {
			qualityCounts[5] = qualityCounts[5] + 1
			if quality >= 6 {
				qualityCounts[6] = qualityCounts[6] + 1
				if quality >= 7 {
					qualityCounts[7] = qualityCounts[7] + 1
				}
			}
		}
	}
	fmt.Printf("Quality over 5:	%v\n", qualityCounts[5])
	fmt.Printf("Quality over 6:	%v\n", qualityCounts[6])
	fmt.Printf("Quality over 7:	%v\n", qualityCounts[7])
	fmt.Printf("Tuner cost:	%v\n", tunerCost)
	//log.Printf("Quality/Tuners:	%v", (qualitySum/float32(totalDiscs))*float32(totalDiscs)/float32(tunerCost))
}

func TestWithChosenSubstatsValue(t *testing.T) {
	subQualityMap := ruptureSubQualityMap
	wantedMainstat := HPP
	WithChosenSubstatsValueSingleTest(subQualityMap, 120000, wantedMainstat, CritDmg, CritRate)
	WithChosenSubstatsValueSingleTest(subQualityMap, 200000, wantedMainstat, CritDmg)
	WithChosenSubstatsValueSingleTest(subQualityMap, 600000, wantedMainstat)
}

// Shared between tests

var attackerSubQualityMap map[stat]float32 = map[stat]float32{
	HP:                 0.0,
	HPP:                0.0,
	DEF:                0.0,
	DEFP:               0.0,
	ATK:                0.3,
	ATKP:               0.9,
	CritRate:           1.0,
	CritDmg:            1.0,
	AnomalyProficiency: 0.1,
	PEN:                0.4,
}

var ruptureSubQualityMap map[stat]float32 = map[stat]float32{
	HP:                 0.3,
	HPP:                0.7,
	DEF:                0.0,
	DEFP:               0.0,
	ATK:                0.1,
	ATKP:               0.3,
	CritRate:           1.0,
	CritDmg:            1.0,
	AnomalyProficiency: 0.1,
	PEN:                0.0,
}

var anomalySubQualityMap map[stat]float32 = map[stat]float32{
	HP:                 0.0,
	HPP:                0.0,
	DEF:                0.0,
	DEFP:               0.0,
	ATK:                0.3,
	ATKP:               0.9,
	CritRate:           0.2,
	CritDmg:            0.1,
	AnomalyProficiency: 1.0,
	PEN:                0.4,
}

var miyabiSubQualityMap map[stat]float32 = map[stat]float32{
	HP:                 0.0,
	HPP:                0.0,
	DEF:                0.0,
	DEFP:               0.0,
	ATK:                0.3,
	ATKP:               0.9,
	CritRate:           1.0,
	CritDmg:            1.0,
	AnomalyProficiency: 0.9,
	PEN:                0.4,
}
