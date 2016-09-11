package main

import "krishnak.org/krishproxy/server"

func main() {
	gameServer := server.NewGameServer()
	gameServer.Start("", 4500)
}
