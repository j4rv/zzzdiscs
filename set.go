package zzzdiscs

type DiscSet string

var SetNameToID = map[DiscSet]int{
	"ShiningAria":       336,
	"WhiteWaterBallad":  335,
	"MoonlightLullaby":  334,
	"DawnsBloom":        333,
	"KingOfTheSummit":   332,
	"YunkuiTales":       331,
	"PhaethonsMelody":   330,
	"ShadowHarmony":     329,
	"AstralVoice":       328,
	"BranchBladeSong":   327,
	"FangedMetal":       326,
	"PolarMetal":        325,
	"ThunderMetal":      324,
	"ChaoticMetal":      323,
	"InfernoMetal":      322,
	"ProtoPunk":         319,
	"ChaosJazz":         318,
	"SwingJazz":         316,
	"SoulRock":          315,
	"HormonePunk":       314,
	"FreedomBlues":      313,
	"ShockstarDisco":    312,
	"PufferElectro":     311,
	"WoodpeckerElectro": 310,
}

var AllDiscSets []DiscSet

func init() {
	AllDiscSets = make([]DiscSet, 0, len(SetNameToID))
	for k := range SetNameToID {
		AllDiscSets = append(AllDiscSets, k)
	}
}
