package traceFall

import (
	"encoding/json"
)

type ExtraData map[string]interface{}

func (e *ExtraData) Set(key string, val interface{}) *ExtraData {
	(*e)[key] = val
	return e
}

func (e ExtraData) Get(key string) interface{} {
	if val, ok := e[key]; ok {
		return val
	}
	return nil
}

func (e ExtraData) ToJson() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		b = []byte(`{}`)
	}
	return b
}

func (e *ExtraData) FromJson(str string) {
	json.Unmarshal([]byte(str), e)
}

func NewExtraData() ExtraData {
	return make(ExtraData)
}
