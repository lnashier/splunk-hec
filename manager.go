package hec

import (
	"context"
	"time"
)

type EventManager struct {
	client         *Client
	events         chan any
	exit           chan struct{}
	batchSize      int
	maxWaitSeconds int
}

func NewEventManager(opt ...ManagerOpt) *EventManager {
	opts := defaultManagerOpts
	opts.apply(opt)

	return &EventManager{
		client:         opts.client,
		batchSize:      opts.batchSize,
		maxWaitSeconds: opts.maxWaitSeconds,
		events:         make(chan any, opts.bufferSize),
		exit:           make(chan struct{}),
	}
}

func (em *EventManager) Start() error {
	go em.Dispatch()
	return nil
}

func (em *EventManager) Stop() error {
	close(em.exit)
	return nil
}

func (em *EventManager) Publish(event any) {
	select {
	case <-em.exit:
	case em.events <- event:
	default:
	}
}

func (em *EventManager) Dispatch() {
	ticker := time.NewTicker(time.Duration(em.maxWaitSeconds) * time.Second)

	defer func() {
		ticker.Stop()
	}()

	processEvents := func(events []any) error {
		_, err := em.client.Send(context.Background(), events...)
		return err
	}

	var eventsBatch []any

	for {
		select {
		case <-em.exit:
			return
		case event, ok := <-em.events:
			if !ok {
				if len(eventsBatch) > 0 {
					processEvents(eventsBatch)
					eventsBatch = nil
				}
				return
			}
			eventsBatch = append(eventsBatch, event)
			if len(eventsBatch) >= em.batchSize {
				if err := processEvents(eventsBatch); err != nil {
					return
				}
				eventsBatch = nil
			}
		case <-ticker.C:
			if len(eventsBatch) > 0 {
				if err := processEvents(eventsBatch); err != nil {
					return
				}
				eventsBatch = nil
			}
		}
	}
}
