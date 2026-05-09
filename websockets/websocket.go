package websockets

import (
	"hash/fnv"
	"log/slog"
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

// WebSocketClient is the server-side client of a websocket connection started by a request
type WebSocketClient struct {
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
		err := c.Conn.Close()
		if err != nil {
			slog.Warn("Could not close connection in defer of handleWriting", "err", err, "clientId", c.websocketClientId.String())
		}
	}()
	for {
		select {
		case event := <-c.ServerSideEvents:
			slog.Debug("event sent to client", "event", event, "clientId", c.websocketClientId.String())
			err := c.Conn.WriteJSON(event)
			if err != nil {
				return
			}
		case <-ticker.C:
			slog.Debug("ping", "clientId", c.websocketClientId.String())
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				slog.Warn("Could not set deadline", "err", err, "clientId", c.websocketClientId.String())
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Warn("error occured while sending ping to client", "err", err, "clientId", c.websocketClientId.String())
				return
			}
		}
	}
}

func (c *WebSocketClient) handleReading(hub *hub.Hub) {
	defer func() {
		hub.Unregister(c.ServerSideEvents)
		err := c.Conn.Close()
		if err != nil {
			slog.Warn("Could not close connection in defer of handleReading", "err", err, "clientId", c.websocketClientId.String())
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		slog.Warn("Could not set SetReadLine", "err", err, "clientId", c.websocketClientId.String())
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		slog.Debug("pong", "clientId", c.websocketClientId.String())
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})

	for {
		typus, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("Closing event occured whilst reading/waiting for message", "err", err, "clientId", c.websocketClientId.String())
				return
			}
			break
		}
		slog.Info("Recieved message", "typus", typus, "content", message, "clientId", c.websocketClientId.String())
	}
}

// RegisterRoutes registers the websocket route to the passed router
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		db := requestcontext.DB(r)
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
				slog.Error("Could not find user in database", "err", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		s := WebSocketClient{
			Conn:              conn,
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
