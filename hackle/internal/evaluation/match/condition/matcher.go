package condition

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type Matcher interface {
	Matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error)
}
