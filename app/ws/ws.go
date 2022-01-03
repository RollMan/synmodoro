package ws

import (
  "log"
  "net/http"
  "time"
  "github.com/gorilla/websocket"
  "encoding/json"
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
    Broadcast:  make(chan []byte, 16),
    register:   make(chan *Client, 16),
    unregister: make(chan *Client, 16),
    clients:    make(map[*Client]bool),
  }
}

func (h *Hub) Run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
      log.Println("====registered====")
      log.Println(string(client.Username))

      response, err := h.CreateMatesJson()
      if err != nil {
        log.Println("Failed to parse json of mates.")
        log.Println(err)
      }else{
        log.Println("broadcasting")
        h.Broadcast <- response
      }

    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }
    case message := <-h.Broadcast:
      log.Println("===message===")
      log.Printf("%T\n", message)
      log.Println(string(message))
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          log.Println("closing client")
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}

func (h *Hub) GetClients() map[*Client]bool{
  return h.clients
}

type Client struct {
  hub *Hub

  // The websocket connection.
  conn *websocket.Conn

  // Username
  Username string

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

func (hub *Hub) CreateMatesJson() ([]byte, error) {
  var err error
  var mates []string = []string{}
  for k, v := range hub.GetClients() {
    log.Println("==creating mates==")
    log.Println(string(k.Username), v)
    if v {
      mates = append(mates, k.Username)
    }
  }

  jsonResponse := map[string]interface{} {"Mates": mates}
  response, err := json.Marshal(jsonResponse)
  if err != nil {
    return nil, err
  }
  return response, nil
}


func WsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }

  username, ok := r.URL.Query()["Username"];
  if !ok {
    w.WriteHeader(http.StatusBadRequest)
    log.Println("Invalid query string for ws handler:")
    for k, v := range r.URL.Query() {
      log.Println(k, v)
    }
    return
  }

  if len(username) != 1 {
    w.WriteHeader(http.StatusBadRequest)
    log.Println("Invalid query string for ws handler: too many usernames " + string(len(username)))
  }

  client := &Client{hub: hub, conn: conn, send: make(chan []byte, 16), Username: username[0]}
  client.hub.register <- client

  // Allow collection of memory referenced by the caller by doing all work in
  // new goroutines.
  go client.writePump()
}
