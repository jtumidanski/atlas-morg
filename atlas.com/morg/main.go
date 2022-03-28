package main

import (
	_map "atlas-morg/character"
	"atlas-morg/kafka"
	"atlas-morg/logger"
	"atlas-morg/monster"
	"atlas-morg/rest"
	"atlas-morg/tracing"
	"atlas-morg/world"
	"context"
	"github.com/opentracing/opentracing-go"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-morg"
const consumerGroupId = "Monster Registry Service"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err = tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	kafka.CreateConsumers(l, ctx, wg,
		_map.MapConsumer(consumerGroupId),
		monster.DamageConsumer(consumerGroupId),
		monster.MovementConsumer(consumerGroupId))

	rest.CreateService(l, ctx, wg, "/ms/morg", monster.InitResource, world.InitResource)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()

	span := opentracing.StartSpan("shutdown")
	monster.DestroyAll(l, span)
	span.Finish()

	l.Infoln("Service shutdown.")
}
