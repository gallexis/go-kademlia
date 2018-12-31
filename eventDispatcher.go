package main

import "time"

type event struct {
    endsAt         time.Time
}

type Dispatcher struct {
    out              <-chan bool
    Tick             <-chan time.Time
    EventDurationMax time.Duration
    Map              map[string]event
}

func NewDispatcher() Dispatcher {
    return Dispatcher{
        out:              make(chan bool),
        Map:              make(map[string]event),
        Tick:             time.Tick(time.Second * 15),
        EventDurationMax: time.Second * 5,
    }
}

func (d *Dispatcher) Start() {

    go func() {

        select {
        case <-d.Tick:
            now := time.Now()

            for k, v := range d.Map {
                if now.After(v.endsAt) {
                    delete(d.Map, k)
                }
            }

        case <-d.out:
            return
        }

    }()
}

func (d *Dispatcher) AddEvent(tx string) {
    d.Map[tx] = event{
        endsAt:         time.Now().Add(d.EventDurationMax),
    }
}

func (d *Dispatcher) EventExists(tx string) bool {
    _, ok := d.Map[tx]
    if ok {
        delete(d.Map, tx)
        return true
    }

    return false
}
