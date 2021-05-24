package websockets

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// pingPeriod = (pongWait * 9) / 10
	pingPeriod = (pongWait * 7) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketClient struct {
	log              *logrus.Logger
	Conn             *websocket.Conn
	User             *models.User
	ServerSideEvents chan hub.ServerSideEvent
}

func (c *WebSocketClient) HandleWriting() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case event := <-c.ServerSideEvents:
			log.Println(event)
			c.Conn.WriteJSON(event)
		case <-ticker.C:
			log.Println("Ping")
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (c *WebSocketClient) HandleReading(hub *hub.Hub) {
	defer func() {
		hub.Unregister(c.ServerSideEvents)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		log.Println("Pong Handler")
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		typus, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		log.Printf("Recieved message of typus %d with the content %s", typus, message)
	}
}

func RegisterRoutes(r *mux.Router, db *gorm.DB, log *logrus.Logger, store sessions.Store, h *hub.Hub) {
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			return
		}
		s := WebSocketClient{
			Conn:             conn,
			log:              log,
			ServerSideEvents: make(chan hub.ServerSideEvent, 256),
		}

		h.Register(s.ServerSideEvents)

		go s.HandleWriting()
		go s.HandleReading(h)
	})
}
