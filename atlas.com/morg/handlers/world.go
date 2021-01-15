package handlers

import (
	"atlas-morg/attributes"
	"atlas-morg/models"
	"atlas-morg/processors"
	"atlas-morg/registries"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// KeyWorld is a key used for the Monster object in the context
type KeyWorld struct{}

type World struct {
	l *log.Logger
}

func NewWorld(l *log.Logger) *World {
	return &World{l}
}

func (w *World) GetMonstersInMap(rw http.ResponseWriter, r *http.Request) {
	ms := registries.GetMonsterRegistry().GetMonstersInMap(getWorldId(r), getChannelId(r), getMapId(r))

	var response attributes.MonsterListDataContainer
	response.Data = make([]attributes.MonsterData, 0)
	for _, x := range ms {
		response.Data = append(response.Data, getMonsterResponseObject(&x))
	}

	err := attributes.ToJSON(response, rw)
	if err != nil {
		w.l.Println("Error encoding GetMonstersInMap response")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (w *World) CreateMonsterInMap(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(KeyWorld{}).(attributes.MonsterInputDataContainer).Data.Attributes
	m, err := processors.NewMonster(w.l).CreateMonster(getWorldId(r), getChannelId(r), getMapId(r), &input)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	var response attributes.MonsterDataContainer
	response.Data = getMonsterResponseObject(m)

	err = attributes.ToJSON(response, rw)
	if err != nil {
		w.l.Println("Error encoding GetMonstersInMap response")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func getMonsterResponseObject(m *models.Monster) attributes.MonsterData {
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

func getWorldId(r *http.Request) byte {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars["worldId"])
	if err != nil {
		log.Println("Error parsing worldId as integer")
		return 0
	}
	return byte(value)
}

func getChannelId(r *http.Request) byte {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars["channelId"])
	if err != nil {
		log.Println("Error parsing channelId as integer")
		return 0
	}
	return byte(value)
}

func getMapId(r *http.Request) int {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars["mapId"])
	if err != nil {
		log.Println("Error parsing mapId as integer")
		return 0
	}
	return value
}
