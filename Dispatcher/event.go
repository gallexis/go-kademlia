package Dispatcher

import "time"

type Event struct {
    startTime  time.Time
    Retries    int
    Duplicates int
    OnTimeout  Callback
    OnResponse Callback
    OnRetry    Callback
}

func (e Event) HasTimedOut(now time.Time, timeout time.Duration) bool  {
    return now.Before(e.startTime.Add(timeout))
}
