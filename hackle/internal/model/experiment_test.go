package model

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExperiment(t *testing.T) {

	experiment := Experiment{}
	_, ok := experiment.WinnerVariation()
	assert.Equal(t, false, ok)

	experiment = Experiment{
		Variations:        []Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
		WinnerVariationID: ref.Int64(1001),
	}
	variation, ok := experiment.WinnerVariation()
	assert.Equal(t, Variation{ID: 1001, Key: "A"}, variation)
	assert.Equal(t, true, ok)

	_, ok = experiment.GetVariationByID(1000)
	assert.Equal(t, false, ok)
	variation, ok = experiment.GetVariationByID(1002)
	assert.Equal(t, Variation{ID: 1002, Key: "B"}, variation)
	assert.Equal(t, true, ok)

	_, ok = experiment.GetVariationByKey("C")
	assert.Equal(t, false, ok)
	variation, ok = experiment.GetVariationByKey("B")
	assert.Equal(t, Variation{ID: 1002, Key: "B"}, variation)
	assert.Equal(t, true, ok)
}

func TestNewExperimentStatusFrom(t *testing.T) {

	test := func(value string, status ExperimentStatus, ok bool) {
		a, b := NewExperimentStatusFrom(value)
		assert.Equal(t, status, a)
		assert.Equal(t, ok, b)
	}

	test("READY", ExperimentStatusDraft, true)
	test("RUNNING", ExperimentStatusRunning, true)
	test("PAUSED", ExperimentStatusPaused, true)
	test("STOPPED", ExperimentStatusCompleted, true)
	test("INVALID", "", false)
}
