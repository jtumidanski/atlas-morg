package requests

import (
	"atlas-morg/json"
	"atlas-morg/rest/attributes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Map struct {
	l logrus.FieldLogger
}

func NewMap(l logrus.FieldLogger) *Map {
	return &Map{l}
}

func (m *Map) GetCharacterIdsInMap(worldId byte, channelId byte, mapId int) ([]int, error) {
	r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/mrg/worlds/%d/channels/%d/maps/%d/characters", worldId, channelId, mapId))
	if err != nil {
		m.l.WithError(err).Errorf("Retrieving information for characters in the map")
		return nil, err
	}

	td := &attributes.MapCharactersListDataContainer{}
	err = json.FromJSON(td, r.Body)
	if err != nil {
		m.l.WithError(err).Errorf("Decoding information for characters in the map")
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
