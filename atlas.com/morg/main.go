package main

import (
	"atlas-morg/kafka/consumers"
	"atlas-morg/logger"
	"atlas-morg/monster"
	"atlas-morg/rest"
	"atlas-morg/world"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	consumers.CreateEventConsumers(l, ctx, wg)

	rest.CreateService(l, ctx, wg, "/ms/morg", monster.InitResource, world.InitResource)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")

	monster.DestroyAll(l)
}
