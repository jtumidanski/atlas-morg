package monster

import (
	"atlas-morg/rest/attributes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetMonster(l *logrus.Logger) http.HandlerFunc {
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
		var response attributes.MonsterDataContainer
		response.Data = getMonsterResponseObject(model)

		err = attributes.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Println("Encoding GetMonster response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}

}
func getMonsterResponseObject(m *Model) attributes.MonsterData {
	var monsterDamage []attributes.MonsterDamage
	for _, x := range m.DamageEntries() {
		monsterDamage = append(monsterDamage, attributes.MonsterDamage{
			CharacterId: x.CharacterId,
			Damage:      x.Damage,
		})
	}

	return attributes.MonsterData{
		Id:   strconv.Itoa(m.UniqueId()),
		Type: "com.atlas.morg.rest.attribute.MonsterAttributes",
		Attributes: attributes.MonsterAttributes{
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
			Hp:                 m.Hp(),
			DamageEntries:      monsterDamage,
		},
	}
}