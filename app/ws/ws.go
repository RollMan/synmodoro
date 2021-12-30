package ws

import (
  "log"
  "net/http"
  "time"
  "github.com/gorilla/websocket"
)

const (
  // Time allowed to write a message to the peer.
  writeWait = 10 * time.Second

  // Time allowed to read the next pong message from the peer.
  pongWait = 60 * time.Second

  // Send pings to peer with this period. Must be less than pongWait.
  pingPeriod = (pongWait * 9) / 10

  // Maximum message size allowed from peer.
  maxMessageSize = 512
)

var (
  newline = []byte{'\n'}
  space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

type Hub struct {
  // Registered clients.
  clients map[*Client]bool

  // Inbound messages from the clients.
  Broadcast chan []byte

  // Register requests from the clients.
  register chan *Client

  // Unregister requests from clients.
  unregister chan *Client
}

func NewHub() *Hub {
  return &Hub{
    Broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
  }
}

func (h *Hub) Run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }
    case message := <-h.Broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}

type Client struct {
  hub *Hub

  // The websocket connection.
  conn *websocket.Conn

  // Buffered channel of outbound messages.
  send chan []byte
}

func (c *Client) writePump() {
  ticker := time.NewTicker(pingPeriod)
  defer func() {
    c.conn.Close()
  }()

  for {
    select {
    case message, ok := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        log.Println("send chan []byte was not ok.")
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }

      w, err := c.conn.NextWriter(websocket.TextMessage)
      if err != nil {
        return
      }
      w.Write(message)

      if err := w.Close(); err != nil {
        log.Println("Failed to close a write channel.")
        return
      }

    case <- ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
        log.Println(err)
        return
      }
    }
  }

}

func WsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }
  client := &Client{hub: hub, conn: conn, send: make(chan []byte)}
  client.hub.register <- client

  // Allow collection of memory referenced by the caller by doing all work in
  // new goroutines.
  go client.writePump()
}
