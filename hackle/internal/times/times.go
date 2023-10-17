package times

import "time"

func Millis(duration float64) float64 {
	return duration / float64(time.Millisecond)
}
