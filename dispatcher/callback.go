package dispatcher

import (
	"fmt"
	"reflect"
)

// stackoverflow.com/questions/52759729/is-there-a-way-to-define-a-function-can-run-any-callback-function-in-golang

type Callback struct {
	fn   reflect.Value
	args []reflect.Value
}

func (c Callback) isSet() bool {
	return !reflect.DeepEqual(c, Callback{})
}

func (c *Callback) AddArgs(vargs ...interface{}) {
	for _, arg := range vargs {
		c.args = append(c.args, reflect.ValueOf(arg))
	}
}

func (c *Callback) Call(args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Call() : ", r)
		}
	}()

	if c.fn.Kind() != reflect.Func {
		return
	}

	vargs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vargs[i] = reflect.ValueOf(arg)
	}

	c.fn.Call(append(c.args, vargs...))
}

func NewCallback(fn interface{}, args ...interface{}) Callback {
	f := reflect.ValueOf(fn)
	if f.Kind() != reflect.Func {
		panic("not a function")
	}

	vargs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vargs[i] = reflect.ValueOf(arg)
	}

	return Callback{
		fn:   f,
		args: vargs,
	}
}
