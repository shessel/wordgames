package server

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    Name string
    conn* websocket.Conn
    toServer chan string
}

func NewClient(name string, conn *websocket.Conn, toServer chan string) *Client {
    client := &Client{ name, conn, toServer }
    go client.run()
    return client
}

func (client *Client) run() {
    for {
        _, message, err := client.conn.ReadMessage()
        if err != nil {
            break;
        }
        client.toServer <- string(message)
    }
}

func (client *Client) SendMessage(message string) (err error) {
    return client.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (client *Client) Disconnect() {
    client.conn.Close()
}
