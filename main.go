package main

import "github.com/shessel/wordgames/server"

func main() {
    server := server.NewServer()
    server.Start()
}
