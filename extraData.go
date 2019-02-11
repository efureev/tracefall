package traceFall

import (
	"encoding/json"
)

// ExtraData is data tree
type ExtraData map[string]interface{}

func (e *ExtraData) Set(key string, val interface{}) *ExtraData {
	(*e)[key] = val
	return e
}

func (e *ExtraData) Clear() *ExtraData {
	*e = NewExtraData()
	return e
}

func (e ExtraData) Get(key string) interface{} {
	if val, ok := e[key]; ok {
		return val
	}
	return nil
}

func (e ExtraData) ToJSON() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		b = []byte(`{}`)
	}
	return b
}

func (e *ExtraData) FromJSON(b []byte) error {
	return json.Unmarshal(b, e)
}

func NewExtraData() ExtraData {
	return make(ExtraData)
}
