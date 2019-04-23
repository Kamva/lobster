package lobster

// Listener is an interface for event listeners.
type Listener interface {
	// Construct initialize listener dependencies.
	Construct() Listener

	// Handle handles the event and do the related processing.
	Handle(event Event, data interface{})
}

// Rollback is a function that runs when any critical error is panicked.
type Rollback func(interface{})

// EventMap is a map of event names and their event listeners
type EventMap map[string]EventListener

// EventListener is a struct containing Listener and RollBack, assigned to an event.
type EventListener struct {
	// 	Listener is list of listeners.
	Listener []Listener

	// RollBack is a function that run when any critical panic occurred.
	RollBack Rollback
}

// Event is concurrency event handler.
type Event interface {
	// Fire runs the event listeners assigned to given event.
	Fire(event string, data interface{}) bool

	// RecoverRoutinePanic recover panics inside routines and push it to lobster
	// error channel.
	RecoverRoutinePanic(caller string, critical bool)
}
