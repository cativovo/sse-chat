package main

import (
	"reflect"
	"testing"
)

func TestProvider(t *testing.T) {
	t.Run("subscribe and publish", func(t *testing.T) {
		p := NewProvider()
		c, cancel := p.Subscribe()

		events := []Event{
			{
				EventType: EventTypeMessage,
				Data:      1,
			},
			{
				EventType: EventTypeMessage,
				Data:      2,
			},
		}

		go func() {
			for _, e := range events {
				p.Publish(e)
			}
			cancel()
		}()

		got := make([]Event, 0, len(events))
		for e := range c {
			got = append(got, e)
		}

		if !reflect.DeepEqual(got, events) {
			t.Errorf("want %v, got %v", events, got)
		}
	})

	t.Run("cancel subscription", func(t *testing.T) {
		p := NewProvider()
		c, cancel := p.Subscribe()

		wantInitLen := 1
		gotInitLen := len(p.subscribers)
		if gotInitLen != wantInitLen {
			t.Errorf("want %d, got %d", gotInitLen, wantInitLen)
		}

		cancel()

		wantC := false
		_, gotC := <-c
		if gotC != wantC {
			t.Errorf("want %v, got %v", wantC, gotC)
		}

		wantLen := 0
		gotLen := len(p.subscribers)
		if gotLen != wantLen {
			t.Errorf("want %d, got %d", wantLen, gotLen)
		}
	})
}
