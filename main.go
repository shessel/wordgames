package main

import "github.com/shessel/wordgames/wordgames"

func main() {
    server := wordgames.NewServer()
    server.Start()
}
