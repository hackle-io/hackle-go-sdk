package evaluator

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
)

type Request interface {
	Key() Key
	Workspace() workspace.Workspace
	User() user.HackleUser
}

type SimpleRequest struct {
	K Key
	W workspace.Workspace
	U user.HackleUser
}

func (r SimpleRequest) Key() Key {
	return r.K
}

func (r SimpleRequest) Workspace() workspace.Workspace {
	return r.W
}

func (r SimpleRequest) User() user.HackleUser {
	return r.U
}
