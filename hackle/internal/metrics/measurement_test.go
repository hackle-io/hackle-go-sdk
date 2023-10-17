package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestField(t *testing.T) {
	assert.Equal(t, "count", FieldCount.String())
	assert.Equal(t, "total", FieldTotal.String())
	assert.Equal(t, "max", FieldMax.String())
	assert.Equal(t, "mean", FieldMean.String())
}
