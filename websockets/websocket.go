package websockets

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/craftamap/hobbit-tracker/hub"
	"github.com/craftamap/hobbit-tracker/middleware/authtocontext"
	"github.com/craftamap/hobbit-tracker/middleware/requestcontext"
	"github.com/craftamap/hobbit-tracker/models"
	"github.com/google/uuid"
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

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func errorln(c *WebSocketClient, msg string) {
	c.log.Errorln(fmt.Sprintf("[%s]", c.getLipglossStyle().Render(c.websocketClientId.String())), msg)
}

func debugln(c *WebSocketClient, msg string) {
	c.log.Debugln(fmt.Sprintf("[%s]", c.getLipglossStyle().Render(c.websocketClientId.String())), msg)
}

// WebSocketClient is the server-side client of a websocket connection started by a request
type WebSocketClient struct {
	log               *logrus.Logger
	Conn              *websocket.Conn
	AuthDetails       authtocontext.AuthDetails
	User              *models.User
	ServerSideEvents  chan hub.ServerSideEvent
	websocketClientId uuid.UUID
}

func (c *WebSocketClient) rangedUUIDHash() int {
	return int(hash(c.websocketClientId.String()) % 231)
}

func (c *WebSocketClient) getLipglossStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(c.rangedUUIDHash())))
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
			debugln(c, fmt.Sprintln("event sent to client - event:", event))
			err := c.Conn.WriteJSON(event)
			if err != nil {
				return
			}
		case <-ticker.C:
			debugln(c, "ping")
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				errorln(c, fmt.Sprintln("error occured while sending ping to client:", err))
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
		debugln(c, "pong")
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		typus, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				errorln(c, fmt.Sprintln("Closing event occured whilst reading/ waiting for message", err))
			}
			break
		}
		c.log.Printf("Recieved message of typus %d with the content %s", typus, message)
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
			Conn:              conn,
			log:               log,
			ServerSideEvents:  make(chan hub.ServerSideEvent, 256),
			AuthDetails:       authDetails,
			User:              &user,
			websocketClientId: uuid.New(),
		}

		h.Register(s.ServerSideEvents)

		go s.handleWriting()
		go s.handleReading(h)
	})
}
