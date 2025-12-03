package zzzdiscs

import (
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
			}
		}
	}
}
