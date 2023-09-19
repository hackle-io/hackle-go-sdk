package condition

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcherFactory(t *testing.T) {

	factory := NewMatcherFactory(map[model.TargetKeyType]Matcher{
		model.TargetKeyTypeUserId: &user.ConditionMatcher{},
	})

	matcher, err := factory.Get(model.TargetKeyTypeUserId)
	assert.Equal(t, nil, err)
	assert.IsType(t, &user.ConditionMatcher{}, matcher)

	matcher, err = factory.Get(model.TargetKeyTypeUserProperty)
	assert.Equal(t, errors.New("unsupported TargetKeyType [USER_PROPERTY]"), err)
}
