package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

func (s WSServer) Start() {
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		clientID := uuid.New().String()
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
				log.Println("ERROR:", err)
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
