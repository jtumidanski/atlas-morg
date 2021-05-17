package world

import (
	"atlas-morg/monster"
	attributes "atlas-morg/rest/attributes"
	"atlas-morg/rest/resource"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetMonstersInMap(l *logrus.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["worldId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		worldId := byte(value)

		vars = mux.Vars(r)
		value, err = strconv.Atoi(vars["channelId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		channelId := byte(value)

		vars = mux.Vars(r)
		value, err = strconv.Atoi(vars["mapId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing mapId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		mapId := value

		ms := monster.GetMonsterRegistry().GetMonstersInMap(worldId, channelId, mapId)

		var response attributes.MonsterListDataContainer
		response.Data = make([]attributes.MonsterData, 0)
		for _, x := range ms {
			response.Data = append(response.Data, getMonsterResponseObject(&x))
		}

		err = attributes.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Encoding GetMonstersInMap response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func CreateMonsterInMap(l *logrus.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cs := &attributes.MonsterInputDataContainer{}
		err := attributes.FromJSON(cs, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Deserializing monster input")
			rw.WriteHeader(http.StatusBadRequest)
			attributes.ToJSON(&resource.GenericError{Message: err.Error()}, rw)
			return
		}

		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["worldId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		worldId := byte(value)

		vars = mux.Vars(r)
		value, err = strconv.Atoi(vars["channelId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		channelId := byte(value)

		vars = mux.Vars(r)
		value, err = strconv.Atoi(vars["mapId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing mapId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		mapId := value

		m, err := monster.Processor(l).CreateMonster(worldId, channelId, mapId, &cs.Data.Attributes)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		var response attributes.MonsterDataContainer
		response.Data = getMonsterResponseObject(m)

		err = attributes.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Encoding GetMonstersInMap response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getMonsterResponseObject(m *monster.Model) attributes.MonsterData {
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
			MaxHp:              m.MaxHp(),
			Hp:                 m.Hp(),
			MaxMp:              m.MaxMp(),
			Mp:                 m.Mp(),
			DamageEntries:      monsterDamage,
		},
	}
}
