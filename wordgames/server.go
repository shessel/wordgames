package wordgames

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

type Server struct {
    clients map[string]*Client
    input chan string
}

func NewServer() Server {
    return Server { make(map[string]*Client), make(chan string) }
}

func (server *Server) Register(client *Client) {
    _, exists := server.clients[client.Name]
    if exists {
        if err := client.SendMessage(client.Name + " is already logged in"); err != nil {
            log.Print("write:", err)
            return
        }
        client.Disconnect()
    } else {
        server.Broadcast(client.Name + " logged in")
        server.clients[client.Name] = client

        if err := client.SendMessage("You are now logged in as " + client.Name); err != nil {
            log.Print("write:", err)
            return
        }
    }
}

func (server *Server) Unregister(client *Client) {
    client.SendMessage("Server shutting down")
    client.Disconnect()
    delete(server.clients, client.Name)
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (server *Server) Broadcast(message string) {
    for _, client := range server.clients {
        if err := client.SendMessage(message); err != nil {
            log.Print("Error sending message to client " + client.Name)
        }
    }
}

func (server *Server) NewConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }

    _, message, err := conn.ReadMessage()
    if err != nil {
        log.Print("read:", err)
        return
    }

    server.Register(NewClient(string(message), conn, server.input))
}

func (server *Server) Start() {
    go server.run()
    http.HandleFunc("/", server.NewConnection)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func (server *Server) run() {
    for {
        server.Broadcast(<- server.input)
    }
}

func (server *Server) Stop() {
    for _, client := range server.clients {
        server.Unregister(client)
    }
}
