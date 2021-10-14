package monster

import (
	json2 "atlas-morg/json"
	"atlas-morg/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	getMonster = "get_monster"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	r := router.PathPrefix("/monsters").Subrouter()
	r.HandleFunc("/{monsterId}", registerGetMonster(l)).Methods(http.MethodGet)
}

func registerGetMonster(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getMonster, func(span opentracing.Span) http.HandlerFunc {
		return parseMonsterId(l, func(monsterId uint32) http.HandlerFunc {
			return handleGetMonster(l)(span)(monsterId)
		})
	})
}

type monsterIdHandler func(monsterId uint32) http.HandlerFunc

func parseMonsterId(l logrus.FieldLogger, next monsterIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		monsterId, err := strconv.Atoi(vars["monsterId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing monsterId as uint32")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(monsterId))(w, r)
	}
}

func handleGetMonster(l logrus.FieldLogger) func(span opentracing.Span) func(monsterId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(monsterId uint32) http.HandlerFunc {
		return func(monsterId uint32) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				model, err := GetMonsterRegistry().GetMonster(monsterId)
				if err != nil {
					rw.WriteHeader(http.StatusNotFound)
					return
				}
				var response DataContainer
				response.Data = getMonsterResponseObject(model)

				err = json2.ToJSON(response, rw)
				if err != nil {
					l.WithError(err).Println("Encoding handleGetMonster response")
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func getMonsterResponseObject(m *Model) DataBody {
	var monsterDamage []DamageEntry
	for _, x := range m.DamageEntries() {
		monsterDamage = append(monsterDamage, DamageEntry{
			CharacterId: x.CharacterId,
			Damage:      x.Damage,
		})
	}

	return DataBody{
		Id:   strconv.Itoa(int(m.UniqueId())),
		Type: "com.atlas.morg.rest.attribute.MonsterAttributes",
		Attributes: Attributes{
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
