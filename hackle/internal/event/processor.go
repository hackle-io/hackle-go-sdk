package event

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"sync"
	"time"
)

type Processor interface {
	Process(event UserEvent)
	Start()
	Close()
}

func NewProcessor(
	capacity int,
	dispatcher Dispatcher,
	dispatchSize int,
	flushScheduler schedule.Scheduler,
	flushInterval time.Duration,
) Processor {
	return &processor{
		queue:          make(chan message, capacity),
		dispatcher:     dispatcher,
		dispatchSize:   dispatchSize,
		flushScheduler: flushScheduler,
		flushInterval:  flushInterval,
		consumingWait:  &sync.WaitGroup{},
		flushingJob:    nil,
		isStarted:      false,
	}
}

type processor struct {
	queue          chan message
	dispatcher     Dispatcher
	dispatchSize   int
	flushScheduler schedule.Scheduler
	flushInterval  time.Duration
	consumingWait  *sync.WaitGroup
	flushingJob    schedule.Job
	isStarted      bool
	mu             sync.Mutex
}

func (p *processor) Process(event UserEvent) {
	select {
	case p.queue <- eventMessage{event}:
		return
	default:
		logger.Info("Event not processed. Exceed event capacity.")
	}
}

func (p *processor) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isStarted {
		logger.Info("EventProcessor is already started.")
		return
	}

	p.consumingWait.Add(1)
	go p.consuming()

	p.flushingJob = p.flushScheduler.SchedulePeriodically(p.flushInterval, p.flushInterval, func() {
		p.queue <- flushMessage{}
	})

	p.isStarted = true
	logger.Info(fmt.Sprintf("EventProcessor started. Flush events every %s", p.flushInterval))
}

func (p *processor) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.isStarted {
		return
	}
	logger.Info("EventProcessor shutting down.")

	if p.flushingJob != nil {
		p.flushingJob.Cancel()
	}

	p.queue <- shutdownMessage{}
	p.consumingWait.Wait()

	p.dispatcher.Close()
	logger.Info("EventProcessor terminated.")
}

func (p *processor) consuming() {
	defer p.consumingWait.Done()
	var events []UserEvent
	for {
		select {
		case msg := <-p.queue:
			switch m := msg.(type) {
			case eventMessage:
				events = append(events, m.event)
				if len(events) >= p.dispatchSize {
					p.dispatch(events)
					events = nil
				}
			case flushMessage:
				p.dispatch(events)
				events = nil
			case shutdownMessage:
				p.dispatch(events)
				return
			}
		}
	}
}

func (p *processor) dispatch(events []UserEvent) {
	if len(events) == 0 {
		return
	}
	p.dispatcher.Dispatch(events)
}

type message interface{}
type eventMessage struct{ event UserEvent }
type flushMessage struct{}
type shutdownMessage struct{}
