package server

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

type Server struct {
    clients map[string]*Client
    input chan Message
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func NewServer() Server {
    return Server { make(map[string]*Client), make(chan Message) }
}

func (server *Server) Start() {
    go server.run()
    http.HandleFunc("/", server.handleConnection)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func (server *Server) Stop() {
    for _, client := range server.clients {
        server.unregisterClient(client)
    }
}

func (server *Server) run() {
    for {
        message := <- server.input
        server.broadcast(message.Client.Name + ": " + message.Text)
    }
}

func (server *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if checkError(err, "upgrade:") {return}

    _, message, err := conn.ReadMessage()
    if checkError(err, "read:") {return}

    client := NewClient(string(message), conn, server.input)
    server.registerClient(client)
    server.runClient(client)
}

func (server *Server) registerClient(client *Client) {
    _, exists := server.clients[client.Name]
    if exists {
        err := client.SendMessage(client.Name + " is already logged in")
        if checkError(err, "write:") {return}
        client.Disconnect()
    } else {
        server.broadcast(client.Name + " logged in")
        server.clients[client.Name] = client
        err := client.SendMessage("You are now logged in as " + client.Name)
        if checkError(err, "write:") {return}
    }
}

func (server *Server) unregisterClient(client *Client) {
    client.SendMessage("Server shutting down")
    client.Disconnect()
    delete(server.clients, client.Name)
}

func (server *Server) runClient(client *Client) {
    client.run()
    server.unregisterClient(client)
}

func (server *Server) broadcast(message string) {
    for _, client := range server.clients {
        checkError(client.SendMessage(message), "Error sending message to client " + client.Name)
    }
}

func checkError(err error, message string) bool {
    ret := err != nil
    if ret {
        log.Print(message, err)
    }
    return ret;
}
