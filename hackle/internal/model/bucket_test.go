package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBucket(t *testing.T) {

	bucket := Bucket{
		Slots: []Slot{
			{0, 10, 1},
		},
	}

	_, ok := bucket.GetSlot(0)
	assert.Equal(t, true, ok)

	_, ok = bucket.GetSlot(9)
	assert.Equal(t, true, ok)

	_, ok = bucket.GetSlot(10)
	assert.Equal(t, false, ok)

	b := Bucket{}
	_, ok = b.GetSlot(10)
	assert.Equal(t, false, ok)
}
