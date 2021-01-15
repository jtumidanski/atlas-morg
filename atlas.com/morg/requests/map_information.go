package requests

import (
   "atlas-morg/attributes"
   "fmt"
   "log"
   "net/http"
)

type MapInformation struct {
   l *log.Logger
}

func NewMapInformation(l *log.Logger) *MapInformation {
   return &MapInformation{l}
}

func (c *MapInformation) GetMonsterInformation(monsterId int) (*attributes.MonsterDataAttributes, error) {
   r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/mis/monsters/%d", monsterId))
   if err != nil {
      c.l.Printf("[ERROR] retrieving monster information for monster %d", monsterId)
      return nil, err
   }

   td := &attributes.MonsterDataDataContainer{}
   err = attributes.FromJSON(td, r.Body)
   if err != nil {
      c.l.Printf("[ERROR] decoding monster information for monster %d", monsterId)
      return nil, err
   }
   return &td.Data.Attributes, nil
}
