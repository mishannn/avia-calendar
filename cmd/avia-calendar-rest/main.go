package main

import (
	"log"

	"github.com/mishannn/avia-calendar/internal/dal/aviasales"
	"github.com/mishannn/avia-calendar/internal/services/tickets"
	"github.com/mishannn/avia-calendar/internal/views/rest"
)

var serverAddress = ":8796"

func main() {
	aviasalesClient, err := aviasales.NewClient("socks5://127.0.0.1:9050")
	if err != nil {
		log.Panicf("can't create aviasales client: %s", err)
	}

	ticketsService := tickets.NewService(aviasalesClient)

	server, err := rest.NewServer(ticketsService)
	if err != nil {
		log.Panicf("can't create server: %s", err)
	}

	err = server.Run(serverAddress)
	if err != nil {
		log.Panicf("can't run server on address '%s': %s", serverAddress, err)
	}
}
