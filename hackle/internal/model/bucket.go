package model

type Bucket struct {
	ID       int64
	Seed     int
	SlotSize int
	Slots    []Slot
}

func (b *Bucket) GetSlot(slotNumber int) (Slot, bool) {
	for _, slot := range b.Slots {
		if slot.contains(slotNumber) {
			return slot, true
		}
	}
	return Slot{}, false
}

type Slot struct {
	StartInclusive int
	EndExclusive   int
	VariationID    int64
}

func (s *Slot) contains(slotNumber int) bool {
	return s.StartInclusive <= slotNumber && slotNumber < s.EndExclusive
}
