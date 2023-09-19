package evaluator

type Context interface {
	Requests() []Request
	Contains(request Request) bool
	AddRequest(request Request)
	RemoveRequest(request Request)
	Evaluations() []Evaluation
	AddEvaluation(evaluation Evaluation)
}

func NewContext() Context {
	return &context{
		requests:    make([]Request, 0),
		evaluations: make([]Evaluation, 0),
	}
}

type context struct {
	requests    []Request
	evaluations []Evaluation
}

func (c *context) Requests() []Request {
	return c.requests
}

func (c *context) Contains(request Request) bool {
	for _, r := range c.requests {
		if r.Key() == request.Key() {
			return true
		}
	}
	return false
}

func (c *context) AddRequest(request Request) {
	c.requests = append(c.requests, request)
}

func (c *context) RemoveRequest(request Request) {
	index := -1
	for i, r := range c.requests {
		if r.Key() == request.Key() {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.requests = append(c.requests[:index], c.requests[index+1:]...)
}

func (c *context) Evaluations() []Evaluation {
	return c.evaluations
}

func (c *context) AddEvaluation(evaluation Evaluation) {
	c.evaluations = append(c.evaluations, evaluation)
}
