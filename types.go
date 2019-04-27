package lobster

import "github.com/Kamva/shark/exceptions"

// Listener is an interface for event listeners.
type Listener interface {
	// Construct initialize listener dependencies.
	Construct() Listener

	// Handle handles the event and do the related processing.
	Handle(event Event, data interface{})
}

// Rollback is a function that runs when any critical error is panicked.
type Rollback func(interface{}, []exceptions.RoutineException)

// EventMap is a map of event names and their event listeners
type EventMap map[string]EventListener

// EventListener is a struct containing Listener and RollBack, assigned to an event.
type EventListener struct {
	// 	Listener is list of listeners.
	Listener []Listener

	// RollBack is a function that run when any critical panic occurred.
	RollBack Rollback
}

// Output is a map of outputs from listener.
type Output map[string]interface{}

// Event is concurrency event handler.
type Event interface {
	// Fire runs the event listeners assigned to given event.
	Fire(event string, data interface{}) (bool, Output)

	// RecoverRoutinePanic recover panics inside routines and push it to lobster
	// error channel.
	RecoverRoutinePanic(caller string, critical bool)

	// AddOutput add an output to the list of listener output.
	AddOutput(caller string, data interface{})
}
