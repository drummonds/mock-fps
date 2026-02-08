package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nibble/mock-fps/internal/models"
	"github.com/nibble/mock-fps/internal/store"
)

type notification struct {
	resourceType string
	resourceID   string
	eventType    string
}

// Dispatcher handles webhook notification delivery.
type Dispatcher struct {
	store  store.Store
	ch     chan notification
	client *http.Client
}

// NewDispatcher creates a new webhook dispatcher.
func NewDispatcher(s store.Store, bufferSize, workers int) *Dispatcher {
	d := &Dispatcher{
		store: s,
		ch:    make(chan notification, bufferSize),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	for i := 0; i < workers; i++ {
		go d.worker()
	}
	return d
}

// Notify enqueues a notification for delivery.
func (d *Dispatcher) Notify(resourceType, resourceID, eventType string) {
	select {
	case d.ch <- notification{resourceType: resourceType, resourceID: resourceID, eventType: eventType}:
	default:
		log.Printf("webhook: notification buffer full, dropping %s %s %s", resourceType, resourceID, eventType)
	}
}

func (d *Dispatcher) worker() {
	for n := range d.ch {
		d.deliver(n)
	}
}

func (d *Dispatcher) deliver(n notification) {
	subs := d.store.MatchSubscriptions(n.resourceType, n.eventType)
	if len(subs) == 0 {
		return
	}

	payload := models.Notification{
		ID:        uuid.New().String(),
		Type:      "notifications",
		Version:   0,
		CreatedOn: time.Now().UTC(),
		Data: models.NotificationData{
			RecordType: n.resourceType,
			EventType:  n.eventType,
			ResourceID: n.resourceID,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("webhook: marshal error: %v", err)
		return
	}

	for _, sub := range subs {
		req, err := http.NewRequest(http.MethodPost, sub.Attributes.CallbackURI, bytes.NewReader(body))
		if err != nil {
			log.Printf("webhook: request build error for %s: %v", sub.Attributes.CallbackURI, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := d.client.Do(req)
		if err != nil {
			log.Printf("webhook: delivery error to %s: %v", sub.Attributes.CallbackURI, err)
			continue
		}
		resp.Body.Close()
	}
}

// Close shuts down the dispatcher.
func (d *Dispatcher) Close() {
	close(d.ch)
}
