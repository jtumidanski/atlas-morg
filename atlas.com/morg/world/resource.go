package world

import (
	json2 "atlas-morg/json"
	"atlas-morg/monster"
	"atlas-morg/rest/resource"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	r := router.PathPrefix("/worlds").Subrouter()

	r.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", ParseMap(l, GetMonstersInMap)).Methods(http.MethodGet)
	r.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", ParseMap(l, CreateMonsterInMap)).Methods(http.MethodPost)
}

type MapHandler func(l logrus.FieldLogger, worldId byte, channelId byte, mapId uint32) http.HandlerFunc

func ParseMap(l logrus.FieldLogger, next MapHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["worldId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as integer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		worldId := byte(value)

		vars = mux.Vars(r)
		value, err = strconv.Atoi(vars["channelId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as integer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		channelId := byte(value)

		vars = mux.Vars(r)
		value, err = strconv.Atoi(vars["mapId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing mapId as integer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mapId := uint32(value)
		next(l, worldId, channelId, mapId)
	}
}

func GetMonstersInMap(l logrus.FieldLogger, worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ms := monster.GetMonsterRegistry().GetMonstersInMap(worldId, channelId, mapId)

		var response monster.DataListContainer
		response.Data = make([]monster.DataBody, 0)
		for _, x := range ms {
			response.Data = append(response.Data, getMonsterResponseObject(x))
		}

		err := json2.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Encoding GetMonstersInMap response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func CreateMonsterInMap(l logrus.FieldLogger, worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cs := &monster.InputDataContainer{}
		err := json2.FromJSON(cs, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Deserializing monster input")
			rw.WriteHeader(http.StatusBadRequest)
			err := json2.ToJSON(&resource.GenericError{Message: err.Error()}, rw)
			if err != nil {
				l.WithError(err).Errorf("Error writing error")
			}
			return
		}

		m, err := monster.CreateMonster(l)(worldId, channelId, mapId, &cs.Data.Attributes)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		var response monster.DataContainer
		response.Data = getMonsterResponseObject(m)

		err = json2.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Encoding GetMonstersInMap response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getMonsterResponseObject(m *monster.Model) monster.DataBody {
	var monsterDamage []monster.DamageEntry
	for _, x := range m.DamageEntries() {
		monsterDamage = append(monsterDamage, monster.DamageEntry{
			CharacterId: x.CharacterId,
			Damage:      x.Damage,
		})
	}

	return monster.DataBody{
		Id:   strconv.Itoa(m.UniqueId()),
		Type: "com.atlas.morg.rest.attribute.MonsterAttributes",
		Attributes: monster.Attributes{
			WorldId:            m.WorldId(),
			ChannelId:          m.ChannelId(),
			MapId:              m.MapId(),
			MonsterId:          m.MonsterId(),
			ControlCharacterId: m.ControlCharacterId(),
			X:                  m.X(),
			Y:                  m.Y(),
			Fh:                 m.Fh(),
			Stance:             m.Stance(),
			Team:               m.Team(),
			MaxHp:              m.MaxHp(),
			Hp:                 m.Hp(),
			MaxMp:              m.MaxMp(),
			Mp:                 m.Mp(),
			DamageEntries:      monsterDamage,
		},
	}
}
