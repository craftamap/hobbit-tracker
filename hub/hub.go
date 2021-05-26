// Package hub is a event hub system managing events by distributing them to subscribed channels
package hub

import "fmt"

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
}

// New creates a new event hub
func New() *Hub {
	return &Hub{
		Subscribers:       make(map[chan ServerSideEvent]bool),
		eventsToBroadcast: make(chan ServerSideEvent, 512),
	}
}

// Run starts the broadcasting of the events
func (h *Hub) Run() {
	// TODO: merge this with New?
	go func() {
		for {
			select {
			case event := <-h.eventsToBroadcast:
				for subscriber, isActive := range h.Subscribers {
					if isActive {
						select {
						case subscriber <- event:
						default:
							// TODO: Better error message
							fmt.Println("Unable to put thing into thing")
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
	fmt.Println("Broadcast event", event)
	fmt.Println(h.Subscribers)
	select {
	case h.eventsToBroadcast <- event:
	default:
		// TODO: Better error message
		fmt.Println("Unable to put thing into thing")
	}
}