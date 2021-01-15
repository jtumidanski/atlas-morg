package morg

import (
   "atlas-morg/consumers"
   "atlas-morg/handlers"
   "context"
   "github.com/go-openapi/runtime/middleware"
   "github.com/gorilla/mux"
   "log"
   "net/http"
   "os"
)

func main() {
   l := log.New(os.Stdout, "morg ", log.LstdFlags|log.Lmicroseconds)

   go consumers.NewMapCharacter(l, context.Background()).Init()
   go consumers.NewMonsterDamage(l, context.Background()).Init()
   go consumers.NewMonsterMovement(l, context.Background()).Init()
   handleRequests(l)
}

func handleRequests(l *log.Logger) {
   //TODO this needs to be updated
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

   l.Fatal(http.ListenAndServe(":8080", router))
}

func commonHeader(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.Header().Add("Content-Type", "application/json")
      next.ServeHTTP(w, r)
   })
}
