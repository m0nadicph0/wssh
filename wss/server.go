package wss

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrades = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type MessageHandler func(clientID string, messageType int, message []byte)

type WSServer struct {
	host              string
	port              int
	handleMessage     MessageHandler
	clientConnections map[string]*websocket.Conn
}

func NewWSServer(host string, port int, fn MessageHandler) *WSServer {
	return &WSServer{
		host:              host,
		port:              port,
		handleMessage:     fn,
		clientConnections: make(map[string]*websocket.Conn),
	}
}

func (s WSServer) Addr() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s WSServer) WriteText(clientID string, message string) error {
	conn, ok := s.clientConnections[clientID]
	if !ok {
		return errors.New("client ID not found")
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (s WSServer) WriteBinary(clientID string, message []byte) error {
	conn, ok := s.clientConnections[clientID]
	if !ok {
		return errors.New("client ID not found")
	}
	return conn.WriteMessage(websocket.BinaryMessage, []byte(message))
}
func (s WSServer) GetClientIDs() []string {
	keys := make([]string, 0, len(s.clientConnections))
	for k := range s.clientConnections {
		keys = append(keys, k)
	}

	return keys
}

func (s WSServer) Start() {
	router := gin.New()
	router.GET("/ws", func(c *gin.Context) {
		clientID := getClientID(c)
		conn, err := upgrades.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		s.clientConnections[clientID] = conn

		for {
			// Read message from browser
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				color.Red(err.Error())
				break
			}
			// Pass the message to the handler
			s.handleMessage(clientID, mt, msg)
		}
	})
	go func() {
		router.Run(s.Addr())
	}()

}

func (s WSServer) BroadcastText(data string) {
	for _, conn := range s.clientConnections {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(data))
	}
}

func (s WSServer) CloseAllClients() {
	for clientID, conn := range s.clientConnections {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		time.Sleep(1 * time.Second)
		conn.Close()
		delete(s.clientConnections, clientID)
	}
}

func (s WSServer) Close(clientID string) error {
	conn, ok := s.clientConnections[clientID]
	if !ok {
		return errors.New("client ID not found")
	}
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(1 * time.Second)
	conn.Close()
	delete(s.clientConnections, clientID)
	return nil
}

func getClientID(c *gin.Context) string {
	clientID, ok := c.GetQuery("cid")
	if !ok {
		return uuid.New().String()
	}
	return clientID
}
