package bucketer

import (
	"bufio"
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

type hasher struct {
	hash int32
}

func (h *hasher) Hash(string, int32) int32 {
	return h.hash
}

func TestNewBucketer(t *testing.T) {
	actual := NewBucketer()
	assert.IsType(t, &bucketer{}, actual)
}

func TestBucketer_Bucketing(t *testing.T) {

	sut := &bucketer{hasher: &hasher{42}}

	_, ok := sut.Bucketing(model.Bucket{Seed: 320, SlotSize: 10000, Slots: make([]model.Slot, 0)}, "42")
	assert.False(t, ok)

	slot := model.Slot{StartInclusive: 0, EndExclusive: 42, VariationID: 420}
	_, ok = sut.Bucketing(model.Bucket{Seed: 320, SlotSize: 10000, Slots: []model.Slot{slot}}, "42")
	assert.False(t, ok)

	slot = model.Slot{StartInclusive: 0, EndExclusive: 43, VariationID: 420}
	actual, ok := sut.Bucketing(model.Bucket{Seed: 320, SlotSize: 10000, Slots: []model.Slot{slot}}, "42")
	assert.True(t, ok)
	assert.Equal(t, slot, actual)

}

func TestBucketer_CalculateSlotNumber_CSV(t *testing.T) {
	sut := bucketer{&murmur3Hasher{}}

	test := func(filename string) {

		readFile, err := os.Open("../../../../testdata/" + filename + ".csv")
		//goland:noinspection GoUnhandledErrorResult
		defer readFile.Close()

		if err != nil {
			fmt.Println(err)
		}
		fs := bufio.NewScanner(readFile)
		fs.Split(bufio.ScanLines)

		for fs.Scan() {
			line := fs.Text()
			row := strings.Split(line, ",")

			seed := i(row[0])
			slotSize := i(row[1])
			value := row[2]
			slotNumber := i(row[3])

			actual := sut.calculateSlotNumber(seed, slotSize, value)
			assert.Equal(t, slotNumber, actual)
		}
	}

	test("bucketing_all")
	test("bucketing_alphabetic")
	test("bucketing_alphanumeric")
	test("bucketing_numeric")
	test("bucketing_uuid")
}

func i(value string) int {
	number, ok := types.AsNumber(value)
	if !ok {
		panic(value)
	}
	return int(number)
}
