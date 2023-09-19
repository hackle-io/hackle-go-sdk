package event

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/http"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"sync"
)

type Dispatcher interface {
	Dispatch(userEvents []UserEvent)
	Close()
}

func NewDispatcher(eventUrl string, httpClient http.Client) Dispatcher {
	return &dispatcher{
		url:        eventUrl + "/api/v2/events",
		httpClient: httpClient,
		wg:         &sync.WaitGroup{},
	}
}

type dispatcher struct {
	url        string
	httpClient http.Client
	wg         *sync.WaitGroup
}

func (d *dispatcher) Dispatch(userEvents []UserEvent) {
	d.wg.Add(1)
	go d.dispatch(userEvents)
}

func (d *dispatcher) dispatch(userEvents []UserEvent) {
	defer d.wg.Done()
	dto := NewPayloadDTO(userEvents)
	err := d.httpClient.PostObj(d.url, dto)
	if err != nil {
		logger.Error("Failed to dispatch events: %v", err)
	}
}

func (d *dispatcher) Close() {
	logger.Info("EventDispatcher shutting down.")
	d.wg.Wait()
	logger.Info("EventDispatcher terminated.")
}
