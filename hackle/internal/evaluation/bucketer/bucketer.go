package bucketer

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type Bucketer interface {
	Bucketing(bucket model.Bucket, identifier string) (model.Slot, bool)
}

func NewBucketer() Bucketer {
	return &bucketer{
		hasher: &murmur3Hasher{},
	}
}

type bucketer struct {
	hasher Hasher
}

func (b *bucketer) Bucketing(bucket model.Bucket, identifier string) (model.Slot, bool) {
	slotNumber := b.calculateSlotNumber(bucket.Seed, bucket.SlotSize, identifier)
	return bucket.GetSlot(slotNumber)
}

func (b *bucketer) calculateSlotNumber(seed int, slotSize int, value string) int {
	hash := b.hasher.Hash(value, int32(seed))
	return int(abs(hash) % int32(slotSize))
}

func abs(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}
