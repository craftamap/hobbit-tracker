package websockets

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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

// WebSocketClient is the server-side client of a websocket connection started by a request
type WebSocketClient struct {
	log              *logrus.Logger
	Conn             *websocket.Conn
	AuthDetails      authtocontext.AuthDetails
	User             *models.User
	ServerSideEvents chan hub.ServerSideEvent
}

func (c *WebSocketClient) handleWriting() {
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

func (c *WebSocketClient) handleReading(hub *hub.Hub) {
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

// RegisterRoutes registers the websocket route to the passed router
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
		log := requestcontext.Log(r)
		h := requestcontext.Hub(r)

		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			return
		}

		authDetails := r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails)
		user := models.User{}

		// TODO: Add error handling here
		if authDetails.Authenticated {
			err = db.Where("ID = ?", r.Context().Value(authtocontext.AuthDetailsContextKey).(authtocontext.AuthDetails).UserID).First(&user).Error
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		s := WebSocketClient{
			Conn:             conn,
			log:              log,
			ServerSideEvents: make(chan hub.ServerSideEvent, 256),
			AuthDetails:      authDetails,
			User:             &user,
		}

		h.Register(s.ServerSideEvents)

		go s.handleWriting()
		go s.handleReading(h)
	})
}
