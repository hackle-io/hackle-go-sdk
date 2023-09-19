package remoteconfig

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEvaluation(t *testing.T) {

	evaluation := NewEvaluationOf(
		decision.ReasonDefaultRule,
		[]evaluator.Evaluation{Evaluation{}},
		model.RemoteConfigParameter{ID: 42},
		ref.Int64(320),
		"return value",
		map[string]interface{}{"a": "b"},
	)
	assert.Equal(t, "DEFAULT_RULE", evaluation.Reason())
	assert.Equal(t, 1, len(evaluation.TargetEvaluations()))
}
