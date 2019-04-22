package lobster

import "github.com/Kamva/shark/exceptions"

// Listener is a function listens on an event.
type Listener func(Event, interface{})

// Rollback is a function that runs when any critical error is panicked.
type Rollback func(interface{})

// EventMap is a map of event names and their event listeners
type EventMap map[string]EventListener

// EventListener is a struct containing Listener and RollBack, assigned to an event.
type EventListener struct {
	Listener []Listener
	RollBack Rollback
}

// Event is concurrency event handler.
type Event interface {
	Fire(event string, data interface{}) []exceptions.RoutineException
	RecoverRoutinePanic(caller string, critical bool)
}
