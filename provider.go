package main

import "sync"

type EventType string

const (
	EventTypeConnect    EventType = "connect"
	EventTypeDisconnect EventType = "disconnect"
	EventTypeMessage    EventType = "message"
)

type Event struct {
	EventType EventType
	Data      any
}

type Provider interface {
	Subscribe() (c chan Event, cancel func())
	Publish(e Event)
}

type provider struct {
	subMu       sync.RWMutex
	subscribers map[string]chan Event
}

func NewProvider() *provider {
	return &provider{
		subMu:       sync.RWMutex{},
		subscribers: make(map[string]chan Event),
	}
}

func (p *provider) Subscribe() (c chan Event, cancel func()) {
	c = make(chan Event)
	id := generateID()

	p.subMu.Lock()
	defer p.subMu.Unlock()

	p.subscribers[id] = c

	return c, func() {
		p.subMu.Lock()
		defer p.subMu.Unlock()

		delete(p.subscribers, id)
		close(c)
	}
}

func (p *provider) Publish(e Event) {
	p.subMu.RLock()
	defer p.subMu.RUnlock()

	for _, s := range p.subscribers {
		s <- e
	}
}
