package remoteconfig

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
)

type Request struct {
	key          evaluator.Key
	workspace    workspace.Workspace
	user         user.HackleUser
	Parameter    model.RemoteConfigParameter
	requiredType types.ValueType
	defaultValue interface{}
}

func NewRequest(
	workspace workspace.Workspace,
	user user.HackleUser,
	parameter model.RemoteConfigParameter,
	requiredType types.ValueType,
	defaultValue interface{},
) Request {
	return Request{
		key:          evaluator.Key{Type: evaluator.TypeRemoteConfig, ID: parameter.ID},
		workspace:    workspace,
		user:         user,
		Parameter:    parameter,
		requiredType: requiredType,
		defaultValue: defaultValue,
	}
}

func (r Request) String() string {
	return "evaluator.Request(type=REMOTE_CONFIG, key=" + r.Parameter.Key + ")"
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
