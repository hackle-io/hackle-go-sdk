package metric

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"strconv"
)

type Operation string

const (
	OperationGetWorkspace Operation = "get.workspace"
	OperationPostEvents   Operation = "post.events"
)

func RecordAPI(operation Operation, sample *metrics.TimerSample, success bool) {
	tags := metrics.Tags{
		"operation": string(operation),
		"success":   strconv.FormatBool(success),
	}
	timer := metrics.NewTimer("api.call", tags)
	sample.Stop(timer)
}
