package _map

import (
	"atlas-morg/rest/response"
	"encoding/json"
)

type characterDataContainer struct {
	data     response.DataSegment
	included response.DataSegment
}

type characterDataBody struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes characterAttributes `json:"attributes"`
}

type characterAttributes struct {
}

func (c *characterDataContainer) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(c.data) == 1 {
		t.Data = c.data[0]
	} else {
		t.Data = c.data
	}
	return json.Marshal(t)
}

func (c *characterDataContainer) UnmarshalJSON(data []byte) error {
	d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyMapCharacterData))
	if err != nil {
		return err
	}
	c.data = d
	return nil
}

func (c *characterDataContainer) Data() *characterDataBody {
	if len(c.data) >= 1 {
		return c.data[0].(*characterDataBody)
	}
	return nil
}

func (c *characterDataContainer) DataList() []characterDataBody {
	var r = make([]characterDataBody, 0)
	for _, x := range c.data {
		r = append(r, *x.(*characterDataBody))
	}
	return r
}

func EmptyMapCharacterData() interface{} {
	return &characterDataBody{}
}
