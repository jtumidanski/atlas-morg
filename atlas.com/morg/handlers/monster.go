package handlers

import (
   "atlas-morg/attributes"
   "atlas-morg/registries"
   "github.com/gorilla/mux"
   "log"
   "net/http"
   "strconv"
)

type Monster struct {
   l *log.Logger
}

func NewMonster(l *log.Logger) *Monster {
   return &Monster{l}
}

func (m *Monster) GetMonster(rw http.ResponseWriter, r *http.Request) {
   model, err := registries.GetMonsterRegistry().GetMonster(getMonsterId(r))
   if err == nil {
      rw.WriteHeader(http.StatusNotFound)
      return
   }
   var response attributes.MonsterDataContainer
   response.Data = getMonsterResponseObject(model)

   err = attributes.ToJSON(response, rw)
   if err != nil {
      m.l.Println("Error encoding GetMonster response")
      rw.WriteHeader(http.StatusInternalServerError)
   }
}

func getMonsterId(r *http.Request) int {
   vars := mux.Vars(r)
   value, err := strconv.Atoi(vars["monsterId"])
   if err != nil {
      log.Println("Error parsing worldId as integer")
      return 0
   }
   return value
}
