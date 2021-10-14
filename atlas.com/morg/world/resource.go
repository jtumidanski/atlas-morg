package world

import (
	json2 "atlas-morg/json"
	"atlas-morg/monster"
	"atlas-morg/rest"
	"atlas-morg/rest/resource"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	getMonstersInMap   = "get_monsters_in_map"
	createMonsterInMap = "create_monster_in_map"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	r := router.PathPrefix("/worlds").Subrouter()
	r.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", registerGetMonstersInMap(l)).Methods(http.MethodGet)
	r.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/monsters", registerCreateMonsterInMap(l)).Methods(http.MethodPost)
}

func registerGetMonstersInMap(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getMonstersInMap, func(span opentracing.Span) http.HandlerFunc {
		return parseMap(l, func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
			return handleGetMonstersInMap(l)(span)(worldId, channelId, mapId)
		})
	})
}

type mapHandler func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc

func parseMap(l logrus.FieldLogger, next mapHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		worldId, err := strconv.ParseUint(vars["worldId"], 10, 8)
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		channelId, err := strconv.ParseUint(vars["channelId"], 10, 8)
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mapId, err := strconv.ParseUint(vars["mapId"], 10, 32)
		if err != nil {
			l.WithError(err).Errorf("Error parsing mapId as uint32")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(byte(worldId), byte(channelId), uint32(mapId))(w, r)
	}
}

func handleGetMonstersInMap(l logrus.FieldLogger) func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
		return func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				ms := monster.GetMonsterRegistry().GetMonstersInMap(worldId, channelId, mapId)

				var response monster.DataListContainer
				response.Data = make([]monster.DataBody, 0)
				for _, x := range ms {
					response.Data = append(response.Data, getMonsterResponseObject(x))
				}

				err := json2.ToJSON(response, rw)
				if err != nil {
					l.WithError(err).Errorf("Encoding handleGetMonstersInMap response")
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func registerCreateMonsterInMap(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(createMonsterInMap, func(span opentracing.Span) http.HandlerFunc {
		return parseMap(l, func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
			return handleCreateMonsterInMap(l)(span)(worldId, channelId, mapId)
		})
	})
}

func handleCreateMonsterInMap(l logrus.FieldLogger) func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
		return func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
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

				m, err := monster.CreateMonster(l, span)(worldId, channelId, mapId, &cs.Data.Attributes)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				var response monster.DataContainer
				response.Data = getMonsterResponseObject(m)

				err = json2.ToJSON(response, rw)
				if err != nil {
					l.WithError(err).Errorf("Encoding handleGetMonstersInMap response")
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
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
		Id:   strconv.Itoa(int(m.UniqueId())),
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
