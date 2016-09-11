package main

import "github.com/untermina/krishproxy/server"

func main() {
	gameServer := server.NewGameServer()
	gameServer.Start("", 4500)
}
