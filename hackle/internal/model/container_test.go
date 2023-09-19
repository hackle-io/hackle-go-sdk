package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer(t *testing.T) {

	c := Container{}
	_, ok := c.GetGroup(42)
	assert.Equal(t, false, ok)

	c = Container{
		Groups: []ContainerGroup{
			{ID: 40},
			{ID: 41},
		},
	}
	_, ok = c.GetGroup(42)
	assert.Equal(t, false, ok)

	c = Container{
		Groups: []ContainerGroup{
			{ID: 40},
			{ID: 41},
			{ID: 42},
		},
	}
	_, ok = c.GetGroup(42)
	assert.Equal(t, true, ok)
}
