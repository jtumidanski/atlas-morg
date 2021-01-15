package main

import (
	"atlas-morg/consumers"
	"atlas-morg/handlers"
	"atlas-morg/processors"
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l := log.New(os.Stdout, "morg ", log.LstdFlags|log.Lmicroseconds)

	go consumers.NewMapCharacter(l, context.Background()).Init()
	go consumers.NewMonsterDamage(l, context.Background()).Init()
	go consumers.NewMonsterMovement(l, context.Background()).Init()

	router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/morg").Subrouter()
	router.Use(commonHeader)
	router.Handle("/docs", middleware.Redoc(middleware.RedocOpts{BasePath: "/ms/morg", SpecURL: "/ms/morg/swagger.yaml"}, nil))
	router.Handle("/swagger.yaml", http.StripPrefix("/ms/morg", http.FileServer(http.Dir("/"))))

	m := handlers.NewMonster(l)
	mRouter := router.PathPrefix("/monsters").Subrouter()
	mRouter.HandleFunc("/{monsterId}", m.GetMonster).Methods("GET")

	w := handlers.NewWorld(l)
	wRouter := router.PathPrefix("/worlds").Subrouter()
	wRouter.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", w.GetMonstersInMap).Methods("GET")
	wRouter.Handle("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", w.MiddlewareValidateMonster(w.CreateMonsterInMap)).Methods("POST")

	s := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Println("Starting server on port 8080")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Println("Got signal:", sig)

	processors.NewMonster(l).DestroyAll()

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
