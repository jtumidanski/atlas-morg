package monster

import (
	json2 "atlas-morg/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	r := router.PathPrefix("/monsters").Subrouter()

	r.HandleFunc("/{monsterId}", GetMonster(l)).Methods(http.MethodGet)
}

func GetMonster(l logrus.FieldLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["monsterId"])
		if err != nil {
			l.WithError(err).Errorf("Parsing worldId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		monsterId := value

		model, err := GetMonsterRegistry().GetMonster(monsterId)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		var response DataContainer
		response.Data = getMonsterResponseObject(model)

		err = json2.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Println("Encoding GetMonster response")
			rw.WriteHeader(http.StatusInternalServerError)
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
		Id:   strconv.Itoa(m.UniqueId()),
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
