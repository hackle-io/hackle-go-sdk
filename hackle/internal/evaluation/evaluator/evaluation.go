package evaluator

type Evaluation interface {
	Reason() string
	TargetEvaluations() []Evaluation
}

type SimpleEvaluation struct {
	R string
	E []Evaluation
}

func (e SimpleEvaluation) Reason() string {
	return e.R
}

func (e SimpleEvaluation) TargetEvaluations() []Evaluation {
	return e.E
}
