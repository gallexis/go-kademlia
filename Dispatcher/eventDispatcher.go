package Dispatcher

import (
    "github.com/murlokswarm/log"
    "sync"
    "time"
)

const DefaultEventTimeout = time.Second * 5

type Dispatcher struct {
    sync.Mutex
    out          chan bool
    Tick         <-chan time.Time
    Map          map[string]Event
    callbackChan chan Callback
}

func NewDispatcher(callbackChan chan Callback) Dispatcher {
    return Dispatcher{
        out:          make(chan bool),
        Tick:         time.Tick(time.Second * 4),
        Map:          make(map[string]Event),
        callbackChan: callbackChan,
    }
}

func (d *Dispatcher) Stop() {
    d.out <- true
}

func (d *Dispatcher) Start() {
    for {
        select {
        case <-d.Tick:
            now := time.Now()

            d.Lock()
            for k, event := range d.Map {
                if !event.HasTimedOut(now, DefaultEventTimeout) {
                    continue
                }

                if event.Retries > 0 {
                    event.Retries -= 1
                    d.callbackChan <- event.OnRetry
                    d.Map[k] = event
                } else {
                    d.callbackChan <- event.OnTimeout
                    delete(d.Map, k)
                }
            }
            d.Unlock()

        case <-d.out:
            return
        }
    }
}

func (d *Dispatcher) AddEvent(tx string, event Event) {
    if event.Retries > 0 && !event.OnRetry.isSet() {
        log.Warn("when Retries is > 0, you must set a OnRetry callback")
    }

    d.Lock()
    event.startTime = time.Now()
    d.Map[tx] = event
    d.Unlock()
}

func (d *Dispatcher) GetCallback(tx string) (Callback, bool) {
    d.Lock()
    defer d.Unlock()

    event, exists := d.Map[tx]
    if exists {
        if event.Duplicates <= 0 {
            delete(d.Map, tx)
        } else {
            event.Duplicates -= 1
            d.Map[tx] = event
        }
        return event.OnResponse, true
    }

    return Callback{}, false
}
