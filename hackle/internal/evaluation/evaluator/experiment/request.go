package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
)

type Request struct {
	key                 evaluator.Key
	workspace           workspace.Workspace
	user                user.HackleUser
	Experiment          model.Experiment
	DefaultVariationKey string
}

func NewRequest(
	workspace workspace.Workspace,
	user user.HackleUser,
	experiment model.Experiment,
	defaultVariationKey string,
) Request {
	return Request{
		key:                 evaluator.Key{Type: evaluator.TypeExperiment, ID: experiment.ID},
		workspace:           workspace,
		user:                user,
		Experiment:          experiment,
		DefaultVariationKey: defaultVariationKey,
	}
}

func NewRequestFrom(request evaluator.Request, experiment model.Experiment) Request {
	return Request{
		key:                 evaluator.Key{Type: evaluator.TypeExperiment, ID: experiment.ID},
		workspace:           request.Workspace(),
		user:                request.User(),
		Experiment:          experiment,
		DefaultVariationKey: "A",
	}
}

func (r Request) String() string {
	return fmt.Sprintf("evaluator.Request(type=%s, key=%d)", r.Experiment.Type, r.Experiment.Key)
}

func (r Request) Key() evaluator.Key {
	return r.key
}

func (r Request) Workspace() workspace.Workspace {
	return r.workspace
}

func (r Request) User() user.HackleUser {
	return r.user
}
