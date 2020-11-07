package gotrix

import (
	"reflect"

	"github.com/chanbakjsd/gotrix/debug"
	"github.com/chanbakjsd/gotrix/event"
)

// Handler is the interface that represents the methods the client needs from the handler.
type Handler interface {
	Handle(cli *Client, event event.Content)
	AddHandler(toCall interface{})
}

// AddHandler adds the handler to the list of handlers.
func (c *Client) AddHandler(function interface{}) {
	c.Handler.AddHandler(function)
}

type defaultHandler struct {
	handlers map[event.Type][]reflect.Value
}

func (d *defaultHandler) Handle(cli *Client, event event.Content) {
	handlers, ok := d.handlers[event.ContentOf()]
	debug.Debug(event.ContentOf())
	if !ok {
		return
	}
	for _, v := range handlers {
		go v.Call([]reflect.Value{reflect.ValueOf(cli), reflect.ValueOf(event)})
	}
}

func (d *defaultHandler) AddHandler(function interface{}) {
	typ := reflect.TypeOf(function)
	val := reflect.ValueOf(function)

	// Check function type.
	if typ.Kind() != reflect.Func {
		debug.Fields(map[string]interface{}{
			"type": typ,
		}).Warn("Non-function passed into AddHandler. Ignoring.")
		return
	}
	//nolint:ignore mnd - 2 is the number of parameters in a handler.
	if typ.NumIn() != 2 {
		debug.Fields(map[string]interface{}{
			"type": typ,
		}).Warn("Invalid handler type! Expected func(*Client, eventType). Ignoring.")
		return
	}
	if typ.In(0) != reflect.TypeOf(&Client{}) {
		debug.Fields(map[string]interface{}{
			"type": typ,
		}).Warn("Invalid handler type! Expected func(*Client, eventType). Ignoring.")
		return
	}
	if !typ.In(1).Implements(reflect.TypeOf([]event.Content{}).Elem()) {
		debug.Fields(map[string]interface{}{
			"type": typ,
		}).Warn("Invalid handler type! Expected func(*Client, eventType). Ignoring.")
		return
	}

	// Get event type
	method, _ := typ.In(1).MethodByName("ContentOf")
	result := method.Func.Call([]reflect.Value{reflect.Zero(typ.In(1))})
	eventType := event.Type(result[0].String())

	// Add it to the list of handlers
	if _, ok := d.handlers[eventType]; !ok {
		d.handlers[eventType] = make([]reflect.Value, 0, 1)
	}
	d.handlers[eventType] = append(d.handlers[eventType], val)
	debug.Fields(map[string]interface{}{
		"eventType": eventType,
		"type":      typ,
	}).Debug("Added handler.")
}
