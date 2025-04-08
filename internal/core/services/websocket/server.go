package websocket

import (
	"encoding/json"
	"log"
	"sync"

	wsDomain "ApiSmart/internal/core/domain/websocket"

	gorillaWs "github.com/gorilla/websocket"
)

// Client representa un cliente WebSocket conectado
type Client struct {
	conn     *gorillaWs.Conn
	server   *Server
	send     chan []byte
	sensorID string
}

// Server representa el servidor WebSocket
type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

// NewServer crea una nueva instancia del servidor WebSocket
func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run inicia el servidor WebSocket
func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mutex.Unlock()

		case message := <-s.broadcast:
			s.mutex.Lock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mutex.Unlock()
		}
	}
}

// BroadcastEvent envía un evento a todos los clientes conectados
func (s *Server) BroadcastEvent(event interface{}) {
	message, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error al marcar evento: %v", err)
		return
	}
	s.broadcast <- message
}

// HandleWebSocket maneja una nueva conexión WebSocket
func (s *Server) HandleWebSocket(conn *gorillaWs.Conn) {
	client := &Client{
		conn:   conn,
		server: s,
		send:   make(chan []byte, 256),
	}

	s.register <- client

	go client.writePump()
	go client.readPump()
}

// writePump maneja la escritura de mensajes al cliente
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(gorillaWs.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(gorillaWs.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// readPump maneja la lectura de mensajes del cliente
func (c *Client) readPump() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if gorillaWs.IsUnexpectedCloseError(err, gorillaWs.CloseGoingAway, gorillaWs.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Procesar mensaje recibido
		var event wsDomain.SensorEvent
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Error al deserializar mensaje: %v", err)
			continue
		}

		// Procesar evento según su tipo
		c.server.handleSensorEvent(event)
	}
}
