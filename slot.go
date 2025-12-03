package zzzdiscs

type DiscSlot int

const (
	Slot1 DiscSlot = iota
	Slot2
	Slot4
	Slot3
	Slot5
	Slot6
)

func (t DiscSlot) String() string {
	switch t {
	case Slot1:
		return "Disc 1"
	case Slot2:
		return "Disc 2"
	case Slot3:
		return "Disc 3"
	case Slot4:
		return "Disc 4"
	case Slot5:
		return "Disc 5"
	case Slot6:
		return "Disc 6"
	}
	return "Unknown"
}
