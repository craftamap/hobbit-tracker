// Package hub is a event hub system managing events by distributing them to subscribed channels
package hub

import (
	"github.com/sirupsen/logrus"
)

// ServerSideEventTypus is an alias for string used for ServerSideEvents
type ServerSideEventTypus string

//  constants of EventTypus
const (
	HobbitCreated  ServerSideEventTypus = "HobbitCreated"
	HobbitDeleted  ServerSideEventTypus = "HobbitDeleted"
	HobbitModified ServerSideEventTypus = "HobbitModified"
	RecordCreated  ServerSideEventTypus = "RecordCreated"
	RecordDeleted  ServerSideEventTypus = "RecordDeleted"
	RecordModified ServerSideEventTypus = "RecordModified"
)

// ServerSideEvent simulates a event happening on the Server
type ServerSideEvent struct {
	Typus        ServerSideEventTypus `json:"typus"`
	OptionalData interface{}          `json:"optional_data,omitempty"`
}

// Hub is the actual event Hub
type Hub struct {
	Subscribers       map[chan ServerSideEvent]bool
	eventsToBroadcast chan ServerSideEvent
	log               *logrus.Logger
}

// New creates a new event hub
func New(log *logrus.Logger) *Hub {
	return &Hub{
		Subscribers:       make(map[chan ServerSideEvent]bool),
		eventsToBroadcast: make(chan ServerSideEvent, 512),
		log:               log,
	}
}

// Run starts the broadcasting of the events
func (h *Hub) Run() {
	go func() {
		for {
			select {
			case event := <-h.eventsToBroadcast:
				for subscriber, isActive := range h.Subscribers {
					if isActive {
						select {
						case subscriber <- event:
						default:
							h.log.Warnf("Unable to put event %+v into subscriber %+v", event, subscriber)
						}
					}
				}
			}
		}
	}()
}

// Register adds a channel to the event hub
func (h *Hub) Register(channelToRegister chan ServerSideEvent) {
	h.Subscribers[channelToRegister] = true
}

// Unregister removes a channel from the event hub
func (h *Hub) Unregister(channelToUnregister chan ServerSideEvent) {
	delete(h.Subscribers, channelToUnregister)
}

// Broadcast broadcasts the given event to all of the subscribers
func (h *Hub) Broadcast(event ServerSideEvent) {
	select {
	case h.eventsToBroadcast <- event:
	default:
		h.log.Errorf("Unable to put event %+v into event hub broadcaster %+v", event, h.eventsToBroadcast)
	}
}
