package main

import (
	"flag"
	"os"
)

func main() {
	args := os.Args[1:]

	serverPort := flag.Int("port", 8000, "the port to run the server on")

	if len(args) >= 1 && args[0] == "server" {
		startServer(*serverPort)
	} else {
		consoleGame()
	}
}
