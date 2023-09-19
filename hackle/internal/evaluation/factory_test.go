package evaluation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEvaluators(t *testing.T) {
	experimentEvaluator, remoteConfigEvaluator := NewEvaluators()
	assert.NotNil(t, experimentEvaluator)
	assert.NotNil(t, remoteConfigEvaluator)
}
