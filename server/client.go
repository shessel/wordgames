package server

import (
    "html"
    "github.com/gorilla/websocket"
)

type Client struct {
    Name string
    conn* websocket.Conn
    toServer chan Message
}

func NewClient(name string, conn *websocket.Conn, toServer chan Message) *Client {
    client := &Client{ html.EscapeString(name), conn, toServer }
    return client
}

func (client *Client) run() {
    for {
        _, message, err := client.conn.ReadMessage()
        if err != nil {
            break;
        }
        client.toServer <- Message{html.EscapeString(string(message)), client}
    }
}

func (client *Client) SendMessage(message string) (err error) {
    return client.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (client *Client) Disconnect() {
    client.conn.Close()
}
