package lobster

import (
	"fmt"
	"sync"

	"github.com/Kamva/shark/exceptions"
	"github.com/Kamva/shark/sentry"
)

// Lobster is concurrency event handler.
type Lobster struct {
	eventMap  EventMap
	waitGroup *sync.WaitGroup
	channel   chan exceptions.RoutineException
}

// Fire runs the event listeners assigned to given event.
func (l *Lobster) Fire(event string, data interface{}) []exceptions.RoutineException {
	for _, function := range l.eventMap[event].Listener {
		l.waitGroup.Add(1)
		go function(l, data)
	}

	l.waitGroup.Wait()
	close(l.channel)

	var errors []exceptions.RoutineException
	var criticalErrors []exceptions.RoutineException
	for exception := range l.channel {
		errors = append(errors, exception)
		if exception.Critical {
			criticalErrors = append(criticalErrors, exception)
		}
	}

	if len(criticalErrors) > 0 {
		if RollBack := l.eventMap[event].RollBack; RollBack != nil {
			RollBack(data)
		}

		return criticalErrors
	}

	if len(errors) > 0 {
		sentry.CaptureRoutineException(errors)
	}

	return nil
}

// RecoverRoutinePanic recover panics inside routines and push it to lobster
// error channel.
func (l *Lobster) RecoverRoutinePanic(caller string, critical bool) {
	defer l.waitGroup.Done()

	if err := recover(); err != nil {
		if err, ok := err.(exceptions.GenericException); ok {
			l.channel <- exceptions.RoutineException{
				Message:     err.GetErrorMessage(),
				RoutineName: caller,
				Critical:    critical,
			}
		} else {
			l.channel <- exceptions.RoutineException{
				Message:     fmt.Sprint(err),
				RoutineName: caller,
				Critical:    critical,
			}
		}
	}
}

// NewEvent instantiate lobster object
func NewEvent(eventMap EventMap) *Lobster {
	return &Lobster{
		eventMap:  eventMap,
		waitGroup: getWaitGroup(),
		channel:   make(chan exceptions.RoutineException, 10),
	}
}

func getWaitGroup() *sync.WaitGroup {
	return &sync.WaitGroup{}
}
