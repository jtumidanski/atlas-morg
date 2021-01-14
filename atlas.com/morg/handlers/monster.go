package handlers

import (
	"log"
	"net/http"
)

type Monster struct {
	l *log.Logger
}

func NewMonster(l *log.Logger) *Monster {
	return &Monster{l}
}

func (m *Monster) GetMonster(rw http.ResponseWriter, r *http.Request) {

}