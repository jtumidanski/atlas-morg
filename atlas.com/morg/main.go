package main

import (
	"atlas-morg/kafka/consumers"
	"atlas-morg/logger"
	"atlas-morg/monster"
	"atlas-morg/rest"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l := logger.CreateLogger()

	consumers.CreateEventConsumers(l)

	rest.CreateRestService(l)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Println("Got signal:", sig)

	monster.Processor(l).DestroyAll()
}
