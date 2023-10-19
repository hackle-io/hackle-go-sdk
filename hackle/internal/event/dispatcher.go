package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/http"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	nethttp "net/http"
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
	go func() {
		defer d.wg.Done()
		err := d.dispatch(userEvents)
		if err != nil {
			logger.Error("Failed to dispatch events: %v", err)
		}
	}()
}

func (d *dispatcher) dispatch(userEvents []UserEvent) error {
	req, err := d.createRequest(userEvents)
	if err != nil {
		return err
	}

	res, err := d.httpClient.Execute(req)
	if err != nil {
		return err
	}

	return d.handleResponse(res)
}

func (d *dispatcher) createRequest(userEvents []UserEvent) (*nethttp.Request, error) {
	dto := NewPayloadDTO(userEvents)
	body, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	req, err := nethttp.NewRequest(nethttp.MethodPost, d.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (d *dispatcher) handleResponse(res *nethttp.Response) error {
	defer func() {
		e := res.Body.Close()
		if e != nil {
			logger.Warn("failed to close response body: %v", e)
		}
	}()

	if !http.IsSuccessful(res) {
		return fmt.Errorf("http status code: %d", res.StatusCode)
	}

	return nil
}

func (d *dispatcher) Close() {
	logger.Info("EventDispatcher shutting down.")
	d.wg.Wait()
	logger.Info("EventDispatcher terminated.")
}
