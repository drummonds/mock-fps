package lifecycle

import (
	"log"
	"sync"
	"time"
)

// StatusUpdater is a closure that updates the status of a resource.
// It returns the new status on success, or an error.
type StatusUpdater func(newStatus string) error

// StatusChangeCallback is called after each status transition.
type StatusChangeCallback func(resourceType, resourceID, newStatus string)

type pendingTransition struct {
	resourceType string
	resourceID   string
	chain        StatusChain
	updater      StatusUpdater
}

// Engine manages async status transitions.
type Engine struct {
	stepDelay time.Duration
	onChange  StatusChangeCallback

	mu      sync.Mutex
	standin bool
	pending []pendingTransition
}

// NewEngine creates a new lifecycle engine.
func NewEngine(stepDelayMs int, onChange StatusChangeCallback) *Engine {
	return &Engine{
		stepDelay: time.Duration(stepDelayMs) * time.Millisecond,
		onChange:  onChange,
	}
}

// SetStandIn enables or disables stand-in mode.
// When disabled, any queued transitions are drained and started.
func (e *Engine) SetStandIn(enabled bool) {
	e.mu.Lock()
	e.standin = enabled
	var drain []pendingTransition
	if !enabled {
		drain = e.pending
		e.pending = nil
	}
	e.mu.Unlock()

	for _, p := range drain {
		e.runTransition(p.resourceType, p.resourceID, p.chain, p.updater)
	}
}

// StandInEnabled returns whether stand-in mode is active.
func (e *Engine) StandInEnabled() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.standin
}

// QueueLength returns the number of transitions waiting in the stand-in queue.
func (e *Engine) QueueLength() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.pending)
}

// StartTransition begins an async status chain for a resource.
// If stand-in mode is active, the transition is queued instead.
func (e *Engine) StartTransition(resourceType, resourceID string, chain StatusChain, updater StatusUpdater) {
	e.mu.Lock()
	if e.standin {
		e.pending = append(e.pending, pendingTransition{
			resourceType: resourceType,
			resourceID:   resourceID,
			chain:        chain,
			updater:      updater,
		})
		e.mu.Unlock()
		return
	}
	e.mu.Unlock()

	e.runTransition(resourceType, resourceID, chain, updater)
}

func (e *Engine) runTransition(resourceType, resourceID string, chain StatusChain, updater StatusUpdater) {
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
