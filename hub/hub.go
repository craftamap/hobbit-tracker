// hub is a event hub system managing events by distributing them to subscribed channels
package hub

import "fmt"

type ServerSideEventTypus string

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
	Typus        ServerSideEventTypus
	OptionalData interface{}
}

type Hub struct {
	Subscribers       map[chan ServerSideEvent]bool
	eventsToBroadcast chan ServerSideEvent
}

func New() *Hub {
	return &Hub{
		Subscribers:       make(map[chan ServerSideEvent]bool),
		eventsToBroadcast: make(chan ServerSideEvent, 512),
	}
}

// TODO: merge this with New?
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
							// TODO: Better error message
							fmt.Println("Unable to put thing into thing")
						}
					}
				}
			}
		}
	}()
}

func (h *Hub) Register(channelToRegister chan ServerSideEvent) {
	h.Subscribers[channelToRegister] = true
}

func (h *Hub) Unregister(channelToUnregister chan ServerSideEvent) {
	delete(h.Subscribers, channelToUnregister)
}

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
