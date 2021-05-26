package rest

import (
	"atlas-morg/monster"
	"atlas-morg/world"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes)
}

func ProduceRoutes(l logrus.FieldLogger) http.Handler {
	router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/morg").Subrouter()
	router.Use(CommonHeader)

	mRouter := router.PathPrefix("/monsters").Subrouter()
	mRouter.HandleFunc("/{monsterId}", monster.GetMonster(l)).Methods("GET")

	wRouter := router.PathPrefix("/worlds").Subrouter()
	wRouter.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", world.GetMonstersInMap(l)).Methods(http.MethodGet)
	wRouter.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", world.CreateMonsterInMap(l)).Methods(http.MethodPost)

	return router
}
