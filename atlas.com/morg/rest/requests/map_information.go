package requests

import (
	"atlas-morg/rest/attributes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MapInformation struct {
	l logrus.FieldLogger
}

func NewMapInformation(l logrus.FieldLogger) *MapInformation {
	return &MapInformation{l}
}

func (c *MapInformation) GetMonsterInformation(monsterId int) (*attributes.MonsterDataAttributes, error) {
	r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/mis/monsters/%d", monsterId))
	if err != nil {
		c.l.WithError(err).Errorf("Retrieving monster information for monster %d", monsterId)
		return nil, err
	}

	td := &attributes.MonsterDataDataContainer{}
	err = attributes.FromJSON(td, r.Body)
	if err != nil {
		c.l.WithError(err).Errorf("Decoding monster information for monster %d", monsterId)
		return nil, err
	}
	return &td.Data.Attributes, nil
}
