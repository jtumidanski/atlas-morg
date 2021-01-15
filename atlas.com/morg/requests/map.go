package requests

import (
	"atlas-morg/attributes"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Map struct {
	l *log.Logger
}

func NewMap(l *log.Logger) *Map {
	return &Map{l}
}

func (m *Map) GetCharacterIdsInMap(worldId byte, channelId byte, mapId int) ([]int, error) {
	r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/mrg/worlds/%d/channels/%d/maps/%d/characters", worldId, channelId, mapId))
	if err != nil {
		m.l.Printf("[ERROR] retrieving information for characters in the map")
		return nil, err
	}

	td := &attributes.MapCharactersListDataContainer{}
	err = attributes.FromJSON(td, r.Body)
	if err != nil {
		m.l.Printf("[ERROR] decoding information for characters in the map")
		return nil, err
	}

	var ids []int
	for _, x := range td.Data {
		id, err := strconv.Atoi(x.Id)
		if err == nil {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
