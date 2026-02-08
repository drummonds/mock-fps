package lifecycle

import (
	"log"
	"time"
)

// StatusUpdater is a closure that updates the status of a resource.
// It returns the new status on success, or an error.
type StatusUpdater func(newStatus string) error

// StatusChangeCallback is called after each status transition.
type StatusChangeCallback func(resourceType, resourceID, newStatus string)

// Engine manages async status transitions.
type Engine struct {
	stepDelay time.Duration
	onChange  StatusChangeCallback
}

// NewEngine creates a new lifecycle engine.
func NewEngine(stepDelayMs int, onChange StatusChangeCallback) *Engine {
	return &Engine{
		stepDelay: time.Duration(stepDelayMs) * time.Millisecond,
		onChange:  onChange,
	}
}

// StartTransition begins an async status chain for a resource.
// The updater closure is responsible for actually persisting the status change.
func (e *Engine) StartTransition(resourceType, resourceID string, chain StatusChain, updater StatusUpdater) {
	go func() {
		// Skip the first status (already set at creation time).
		for i := 1; i < len(chain); i++ {
			time.Sleep(e.stepDelay)
			newStatus := chain[i]
			if err := updater(newStatus); err != nil {
				log.Printf("lifecycle: failed to update %s %s to %s: %v", resourceType, resourceID, newStatus, err)
				return
			}
			if e.onChange != nil {
				e.onChange(resourceType, resourceID, newStatus)
			}
		}
	}()
}
