package revel

import (
	"code.google.com/p/go.net/websocket"
	"reflect"
)

var (
	controllerType    = reflect.TypeOf(Controller{})
	controllerPtrType = reflect.TypeOf(&Controller{})
	websocketType     = reflect.TypeOf((*websocket.Conn)(nil))
)

var ActionInvoker actionInvoker

type actionInvoker struct{}

// Instantiate, bind params, and invoke the given action.
func (f actionInvoker) Call(c *Controller, _ FilterChain) {
	// Instantiate the method.
	methodValue := reflect.ValueOf(c.AppController).MethodByName(c.MethodType.Name)

	// Collect the values for the method's arguments.
	var methodArgs []reflect.Value
	for _, arg := range c.MethodType.Args {
		// If they accept a websocket connection, treat that arg specially.
		var boundArg reflect.Value
		if arg.Type == websocketType {
			boundArg = reflect.ValueOf(c.Request.Websocket)
		} else {
			TRACE.Println("Binding:", arg.Name, "as", arg.Type)
			boundArg = c.Params.Bind(arg.Name, arg.Type)
		}
		methodArgs = append(methodArgs, boundArg)
	}

	var resultValue reflect.Value
	if methodValue.Type().IsVariadic() {
		resultValue = methodValue.CallSlice(methodArgs)[0]
	} else {
		resultValue = methodValue.Call(methodArgs)[0]
	}
	if resultValue.Kind() == reflect.Interface && !resultValue.IsNil() {
		c.Result = resultValue.Interface().(Result)
	}
}