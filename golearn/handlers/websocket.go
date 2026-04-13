package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"golearn/config"
	"golearn/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// WSMessage is the message structure exchanged in the classroom.
type WSMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Type     string `json:"type"`     // "message" | "join" | "leave"
	CourseID string `json:"course_id"`
}

// client represents a single WebSocket connection inside a room.
type wsClient struct {
	conn     *websocket.Conn
	send     chan WSMessage
	username string
	courseID string
}

// room holds all clients connected to a specific course classroom.
type room struct {
	mu      sync.Mutex
	clients map[*wsClient]struct{}
}

// hub manages all active rooms.
type hub struct {
	mu    sync.RWMutex
	rooms map[string]*room
}

var globalHub = &hub{rooms: make(map[string]*room)}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *hub) getOrCreateRoom(courseID string) *room {
	h.mu.Lock()
	defer h.mu.Unlock()
	if r, ok := h.rooms[courseID]; ok {
		return r
	}
	r := &room{clients: make(map[*wsClient]struct{})}
	h.rooms[courseID] = r
	return r
}

func (r *room) broadcast(msg WSMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.clients {
		select {
		case c.send <- msg:
		default:
			// Drop message if client channel is full
		}
	}
}

func (r *room) add(c *wsClient) {
	r.mu.Lock()
	r.clients[c] = struct{}{}
	r.mu.Unlock()
}

func (r *room) remove(c *wsClient) {
	r.mu.Lock()
	delete(r.clients, c)
	r.mu.Unlock()
}

// writePump forwards outbound messages from the send channel to the WebSocket.
func writePump(c *wsClient) {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			log.Printf("ws write error: %v", err)
			return
		}
	}
}

// ClassroomWS godoc
// @Summary      WebSocket classroom chat
// @Description  Connect with ?token=<JWT>. Each courseId is an isolated room.
// @Tags         websocket
// @Param        courseId path   string true "Course ID"
// @Param        token    query  string true "JWT token"
// @Router       /ws/classroom/{courseId} [get]
func ClassroomWS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("courseId")

		// Authenticate via query param (browsers can't set headers for WS)
		tokenStr := c.Query("token")
		if tokenStr == "" {
			// Also accept Bearer header for non-browser clients
			hdr := c.GetHeader("Authorization")
			tokenStr = strings.TrimPrefix(hdr, "Bearer ")
		}

		claims := &middleware.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Upgrade HTTP to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("ws upgrade error: %v", err)
			return
		}

		client := &wsClient{
			conn:     conn,
			send:     make(chan WSMessage, 32),
			username: fmt.Sprintf("user-%d", claims.UserID),
			courseID: courseID,
		}

		rm := globalHub.getOrCreateRoom(courseID)
		rm.add(client)

		// Announce join
		rm.broadcast(WSMessage{
			Username: client.username,
			Text:     client.username + " joined the classroom",
			Type:     "join",
			CourseID: courseID,
		})

		go writePump(client)

		// Read loop — receive messages from this client and broadcast them
		defer func() {
			rm.remove(client)
			close(client.send)
			rm.broadcast(WSMessage{
				Username: client.username,
				Text:     client.username + " left the classroom",
				Type:     "leave",
				CourseID: courseID,
			})
		}()

		for {
			var msg WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				break
			}
			msg.Username = client.username
			msg.Type = "message"
			msg.CourseID = courseID
			rm.broadcast(msg)
		}
	}
}
