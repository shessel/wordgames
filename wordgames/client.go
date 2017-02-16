package wordgames

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    Name string
    conn* websocket.Conn
}

func (client *Client) SendMessage(message string) (err error) {
    return client.conn.WriteMessage(websocket.TextMessage, []byte(message))
}
