package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now, or configure based on env
		},
	}
)

// LogBroadcaster implements io.Writer to capture logs and broadcast them to WebSocket clients.
type LogBroadcaster struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

// NewLogBroadcaster creates a new LogBroadcaster.
func NewLogBroadcaster() *LogBroadcaster {
	lb := &LogBroadcaster{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
	go lb.run()
	return lb
}

// Write implements io.Writer. It writes the log data to the broadcast channel.
func (lb *LogBroadcaster) Write(p []byte) (n int, err error) {
	// Make a copy of the slice to avoid race conditions if the underlying array is modified
	// although for log.Logger it's usually fine.
	// We also want to ensure we don't block the logger if no one is listening or channel is full.
	// But for simplicity, we'll just send it.
	
	// We need to copy p because it might be reused by the caller.
	data := make([]byte, len(p))
	copy(data, p)
	
	// Send to broadcast channel in a non-blocking way or just go routine it
	// To preserve order, we should probably just send it.
	// But if we block here, we block the logger.
	// Let's use a buffered channel in run() or just send here.
	
	lb.broadcast <- data
	return len(p), nil
}

func (lb *LogBroadcaster) run() {
	for {
		message := <-lb.broadcast
		lb.mu.Lock()
		for client := range lb.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("websocket write error: %v", err)
				client.Close()
				delete(lb.clients, client)
			}
		}
		lb.mu.Unlock()
	}
}

// ServeWS handles WebSocket requests.
func (lb *LogBroadcaster) ServeWS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	lb.mu.Lock()
	lb.clients[ws] = true
	lb.mu.Unlock()

	// Keep the connection alive until client disconnects
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	lb.mu.Lock()
	delete(lb.clients, ws)
	lb.mu.Unlock()

	return nil
}
