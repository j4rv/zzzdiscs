package zzzdiscs

import (
	"fmt"
	"strings"
)

const indent = "    "

type RRDisc struct {
	Id            int
	Level         int
	Exp           int
	Star          int
	Lock          bool
	Properties    RRStat
	SubProperties []RRStat
}

type RRStat struct {
	Key       int
	BaseValue int
	AddValue  int
}

func (r RRStat) String() string {
	return fmt.Sprintf(
		"%s.{\n%s    .key = %d,\n%s    .base_value = %d,\n%s    .add_value = %d,\n%s}",
		indent, indent, r.Key,
		indent, r.BaseValue,
		indent, r.AddValue,
		indent,
	)
}

func DiscToRRDisc(d Disc) RRDisc {
	// Constant stuff
	disc := RRDisc{
		Level: 15,
		Exp:   0,
		Star:  1,
		Lock:  true,
	}

	// First 3 digits: Set ID, 4th digit: 0-indexed rarity, 5th digit: 1-indexed slot
	disc.Id = SetNameToID[d.Set]*100 + 40 + int(d.Slot) + 1

	// Stats
	disc.Properties = SubstatToRRSubstat(d.MainStat)
	for _, s := range d.SubStats {
		disc.SubProperties = append(disc.SubProperties, SubstatToRRSubstat(*s))
	}

	return disc
}

func SubstatToRRSubstat(s DiscStat) RRStat {
	return RRStat{
		Key:       StatToID[s.Stat],
		BaseValue: s.BaseValue,
		AddValue:  s.Rolls,
	}
}

func FormatRRDisc(d RRDisc) string {
	var sb strings.Builder

	sb.WriteString(".{\n")
	sb.WriteString(fmt.Sprintf("%s.id = %d,\n", indent, d.Id))
	sb.WriteString(fmt.Sprintf("%s.level = %d,\n", indent, d.Level))
	sb.WriteString(fmt.Sprintf("%s.exp = %d,\n", indent, d.Exp))
	sb.WriteString(fmt.Sprintf("%s.star = %d,\n", indent, d.Star))
	sb.WriteString(fmt.Sprintf("%s.lock = %v,\n", indent, d.Lock))

	// Properties
	sb.WriteString(fmt.Sprintf("%s.properties = .{%s},\n", indent, d.Properties.String()))

	// SubProperties
	sb.WriteString(fmt.Sprintf("%s.sub_properties = .{\n", indent))
	for i, sp := range d.SubProperties {
		sb.WriteString(sp.String())
		if i != len(d.SubProperties)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(fmt.Sprintf("%s},\n", indent))
	sb.WriteString("}")

	return sb.String()
}
