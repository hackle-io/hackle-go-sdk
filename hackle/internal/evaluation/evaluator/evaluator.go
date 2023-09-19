package evaluator

type Evaluator interface {
	Evaluate(request Request, context Context) (Evaluation, error)
}

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	TypeExperiment   Type = "EXPERIMENT"
	TypeRemoteConfig Type = "REMOTE_CONFIG"
)

type Key struct {
	Type Type
	ID   int64
}
