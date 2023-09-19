package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSdk(t *testing.T) {
	sdk := NewSdk("abc")
	assert.Equal(t, "abc", sdk.Key)
	assert.Equal(t, "go-sdk", sdk.Name)
}
