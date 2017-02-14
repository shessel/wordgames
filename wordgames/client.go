package wordgames

import "github.com/gorilla/websocket"

type Client struct {
    Name string
    Conn websocket.Conn
}
