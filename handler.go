package gotrix

import (
	"fmt"
	"reflect"

	"github.com/chanbakjsd/gotrix/debug"
	"github.com/chanbakjsd/gotrix/event"
)

// Handler is the interface that represents the methods the client needs from the handler.
type Handler interface {
	Handle(cli *Client, event event.Event)
	AddHandler(toCall interface{}) error
}

// AddHandler adds the handler to the list of handlers.
func (c *Client) AddHandler(function interface{}) error {
	return c.Handler.AddHandler(function)
}

type defaultHandler struct {
	handlers map[event.Type][]reflect.Value
}

func (d *defaultHandler) Handle(cli *Client, event event.Event) {
	handlers, ok := d.handlers[event.Type()]
	debug.Debug("new event: " + event.Type())
	if !ok {
		return
	}
	for _, v := range handlers {
		go v.Call([]reflect.Value{reflect.ValueOf(cli), reflect.ValueOf(event)})
	}
}

func (d *defaultHandler) AddHandler(function interface{}) error {
	typ := reflect.TypeOf(function)
	val := reflect.ValueOf(function)

	// Check function type.
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("AddHandler: expected func(*Client, EventType), got %T instead", function)
	}
	//nolint:gomnd // 2 is the number of parameters in a handler.
	if typ.NumIn() != 2 {
		return fmt.Errorf("AddHandler: expected func(*Client, EventType), got %T instead", function)
	}
	if typ.In(0) != reflect.TypeOf(&Client{}) {
		return fmt.Errorf("AddHandler: expected func(*Client, EventType), got %T instead", function)
	}

	contentInterface := reflect.Zero(typ.In(1)).Interface()
	content, ok := contentInterface.(event.Event)
	if !ok {
		return fmt.Errorf(
			"AddHandler: invalid function input, expected function to take event, takes %T instead",
			contentInterface,
		)
	}

	// Get event type
	eventType := content.Type()

	// Add it to the list of handlers
	if _, ok := d.handlers[eventType]; !ok {
		d.handlers[eventType] = make([]reflect.Value, 0, 1)
	}
	d.handlers[eventType] = append(d.handlers[eventType], val)
	debug.Debug("added handler: event=" + eventType)
	return nil
}
