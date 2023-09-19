package bucketer

import (
	"unsafe"
)

type Hasher interface {
	Hash(data string, seed int32) int32
}

type murmur3Hasher struct{}

func (m *murmur3Hasher) Hash(data string, seed int32) int32 {
	return int32(hash([]byte(data), uint32(seed)))
}

const (
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
)

func hash(data []byte, seed uint32) uint32 {
	h1 := seed

	n := len(data) / 4
	for b := 0; b < n; b++ {
		k1 := *(*uint32)(unsafe.Pointer(&data[b*4]))

		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19)
		h1 = h1*4 + h1 + 0xe6546b64
	}

	tail := data[n*4:]

	var k1 uint32
	switch len(tail) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}
