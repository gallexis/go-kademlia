package main

import (
    "fmt"
    "reflect"
    "sync"
    "time"
)

const DefaultEventTimeout = time.Second * 12

type Callback struct {
    fn   reflect.Value
    args []reflect.Value
}

func (c *Callback) Call() {
    if c.fn.Kind() != reflect.Func {
        return
    }
    c.fn.Call(c.args)
}

func (c *Callback) CallWithArgs(args ...interface{}) {
    vargs := make([]reflect.Value, len(args))
    for i, arg := range args {
        vargs[i] = reflect.ValueOf(arg)
    }

    c.fn.Call(vargs)
}

// stackoverflow.com/questions/52759729/is-there-a-way-to-define-a-function-can-run-any-callback-function-in-golang
func NewCallback(fn interface{}, args ...interface{}) Callback {
    f := reflect.ValueOf(fn)
    if f.Kind() != reflect.Func {
        panic("not a function")
    }

    vargs := make([]reflect.Value, len(args))
    for i, arg := range args {
        vargs[i] = reflect.ValueOf(arg)
    }

    return Callback{f, vargs}
}

type Event struct {
    timeout           time.Time
    maxTries          int
    duplicates        int
    CallbackOnTimeout Callback
    Callback          Callback
    Caller            Callback
}

type Dispatcher struct {
    sync.Mutex
    out  chan bool
    Tick <-chan time.Time
    Map  map[string]Event
}

func NewDispatcher() Dispatcher {
    return Dispatcher{
        out:   make(chan bool),
        Tick:  time.Tick(time.Second * 5),
        Map:   make(map[string]Event),
    }
}

func (d *Dispatcher) Stop() {
    d.out <- true
}

func (d *Dispatcher) Start() {
    go func() {
        for {

            select {
            case <-d.Tick:
                now := time.Now()

                fmt.Println("-->>>>", len(d.Map))
                for k, v := range d.Map {
                    if now.Before(v.timeout.Add(DefaultEventTimeout)) {
                        continue
                    }

                    if v.maxTries <= 1 {
                        fmt.Println("delete event (Start)", v.Callback.fn.String())
                        if !reflect.DeepEqual(v.CallbackOnTimeout, Callback{}) {
                            fmt.Println("reflect.DeepEqual")
                            v.CallbackOnTimeout.Call()
                        }
                        delete(d.Map, k)
                    } else {
                        fmt.Println("Call Callback", v.Callback.fn.String())
                        v.Caller.Call()
                        v.maxTries -= 1
                        d.Map[k] = v
                        fmt.Println("Call Callback - end")
                    }
                }

            case <-d.out:
                return
            }
        }
    }()
}

func (d *Dispatcher) AddEvent(tx string, event Event) {
    d.Map[tx] = event
}

func (d *Dispatcher) GetEvent(tx string) (Callback, bool) {
    v, ok := d.Map[tx]
    if ok {
        if v.duplicates <= 0 {
            fmt.Println("delete event (GetEvent)", v.Callback.fn.String())
            delete(d.Map, tx)
        } else {
            fmt.Println("duplicates-1 : ", v.duplicates)
            v.duplicates -= 1
            d.Map[tx] = v
        }
        return v.Callback, true
    }

    return Callback{}, false
}
