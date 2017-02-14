package wordgames

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

type Server struct {
    clients map[string]*Client
}

func (server *Server) Register(client *Client) {
    server.clients[client.Name] = client
}

func (server *Server) Unregister(client *Client) {
    delete(server.clients, client.Name)
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (server *Server) NewConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer conn.Close()
    messageType, message, err := conn.ReadMessage()
    if err != nil {
        log.Print("read:", err)
        return
    }
    if err = conn.WriteMessage(messageType, message); err != nil {
        log.Print("write:", err)
        return
    }
}

func (server *Server) Start() {
    http.HandleFunc("/", server.NewConnection)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
